package httpserver

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Profiler struct {
	Data map[string]*Stats
	sync.RWMutex
}

type Stat struct {
	Funcname   string   `json:"funcname"`
	Time       float64  `json:"time"`
	Files      []string `json:"files"`
	CallStacks []string `json:"call_stacks"`
	Lines      []string `json:"lines"`
}

type Stats struct {
	Files      []string
	Lines      []string
	Funcname   string
	CallStacks []string
	Count      int64
	Time       *Time
}

func NewProfiler() *Profiler {
	return &Profiler{
		Data: make(map[string]*Stats, 0),
	}
}

var profiler *Profiler

func init() {
	profiler = NewProfiler()
}

func putStatsHandler(w http.ResponseWriter, r *http.Request) {
	profiler.Lock()
	defer profiler.Unlock()

	clen, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := make([]byte, clen)
	_, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var stat *Stat
	err = json.Unmarshal(body[:clen], &stat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f := concatCallStacks(stat.CallStacks)
	if profiler.Data[f] == nil {
		profiler.Data[f] = &Stats{
			Funcname:   stat.Funcname,
			Files:      stat.Files,
			Lines:      stat.Lines,
			CallStacks: stat.CallStacks,
			Time:       NewTime(),
		}

	}

	profiler.Data[f].Time.Set(stat.Time)
	profiler.Data[f].Count++
}

func concatCallStacks(f []string) string {
	return strings.Join(f, ">")
}

func csvRow(s *Stats, callStack string) []string {
	return []string{
		fmt.Sprint(s.Count),
		s.Funcname,
		callStack,
		s.Time.SMax(),
		s.Time.SMin(),
		s.Time.SSum(),
		s.Time.SAvg(s.Count),
		s.Time.SP50(int(s.Count)),
		s.Time.SP99(int(s.Count)),
	}
}

func getStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	rlen := len(profiler.Data) + 1
	rows := make([][]string, rlen)
	rows[0] = []string{"count", "func", "call_stack", "max", "min", "sum", "avg", "p50", "p99"}
	i := 1
	for k, v := range profiler.Data {
		rows[i] = csvRow(v, k)
		i++
	}

	cw := csv.NewWriter(w)
	cw.WriteAll(rows)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	profiler = NewProfiler()
}

func ListenAndServe(hosts ...string) error {
	http.HandleFunc("/put_stats", putStatsHandler)
	http.HandleFunc("/stats", getStatsHandler)
	http.HandleFunc("/reset", resetHandler)

	host := "localhost:58080"
	if len(hosts) > 0 {
		host = hosts[0]
	}

	return http.ListenAndServe(host, nil)
}

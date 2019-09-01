package punc

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tkuchiki/punc/helper"
)

type Client struct {
	httpclient *http.Client
	host       string
}

var client *Client

func init() {
	tr := &http.Transport{
		MaxIdleConns:    20,
		IdleConnTimeout: 120 * time.Second,
	}

	host := "localhost:58080"
	if os.Getenv("PUNC_HOST") != "" {
		host = os.Getenv("PUNC_HOST")
	}

	client = &Client{
		httpclient: &http.Client{Transport: tr},
		host:       host,
	}
}

func (c *Client) PutStats(stats []byte) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/put_stats", c.host), bytes.NewBuffer(stats))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprint(len(stats)))

	_, err = c.httpclient.Do(req)

	return err
}

func Do() int64 {
	if helper.IsDisabled() {
		return 0
	}

	return time.Now().UnixNano()
}

func Done(start int64) {
	if helper.IsDisabled() {
		return
	}

	elapsed := sec(time.Now().UnixNano() - start)

	pc := make([]uintptr, 15)
	n := runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc[:n])
	goroot := runtime.GOROOT()
	files := make([]string, 0)
	callStacks := make([]string, 0)
	lines := make([]string, 0)
	//fs := make([]runtime.Frame, 0)
	var funcname string

	for {
		frame, b := frames.Next()

		if b {
			funcname = filepath.Base(frame.Function)
			if strings.HasPrefix(frame.File, goroot) || strings.HasPrefix(funcname, "punc.") {
				continue
			} else {
				//fs = append(fs, frame)
				files = append([]string{frame.File}, files[0:]...)
				callStacks = append([]string{funcname}, callStacks[0:]...)
				lines = append([]string{fmt.Sprint(frame.Line)}, lines[0:]...)
			}
		}

		if !b {
			break
		}
	}

	names := strings.Split(callStacks[len(callStacks)-1], ".")
	stats := statsToJsonBytes(names[len(names)-1], elapsed, files, callStacks, lines)
	client.PutStats(stats)
}

func sec(i int64) float64 {
	return float64(i) / float64(time.Second)
}

func strSlice(s []string) string {
	return fmt.Sprintf(`"%s"`, strings.Join(s, `", "`))
}

func statsToJsonBytes(funcname string, elapsed float64, files, callStacks, lines []string) []byte {
	s := fmt.Sprintf(`{"funcname": "%s", "time": %f, "files": [%s], "call_stacks": [%s], "lines": [%s]}`,
		funcname, elapsed, strSlice(files), strSlice(callStacks), strSlice(lines))

	return []byte(s)
}

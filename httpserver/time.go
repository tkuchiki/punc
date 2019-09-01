package httpserver

import "fmt"

type Time struct {
	max         float64
	min         float64
	avg         float64
	sum         float64
	p50         float64
	p99         float64
	percentiles []float64
}

func NewTime() *Time {
	return &Time{}
}

func (t *Time) Set(elapsed float64) {
	if t.max < elapsed {
		t.max = elapsed
	}

	if t.min >= elapsed || t.min == 0.0 {
		t.min = elapsed
	}

	t.sum += elapsed

	t.percentiles = append(t.percentiles, elapsed)
}

func (t *Time) Max() float64 {
	return t.max
}

func (t *Time) Min() float64 {
	return t.min
}

func (t *Time) Avg(count int64) float64 {
	return t.sum / float64(count)
}

func (t *Time) Sum() float64 {
	return t.sum
}

func (t *Time) P50(count int) float64 {
	plen := percentRank(count, 50)
	return t.percentiles[plen]
}

func (t *Time) P99(count int) float64 {

	plen := percentRank(count, 99)
	return t.percentiles[plen]
}

func percentRank(l int, n int) int {
	pLen := (l * n / 100) - 1
	if pLen < 0 {
		pLen = 0
	}

	return pLen
}

func round(f float64) string {
	return fmt.Sprintf("%.6f", f)
}

func (t *Time) SMax() string {
	return round(t.max)
}

func (t *Time) SMin() string {
	return round(t.min)
}

func (t *Time) SAvg(count int64) string {
	return round(t.Avg(count))
}

func (t *Time) SSum() string {
	return round(t.sum)
}

func (t *Time) SP50(count int) string {
	return round(t.P50(count))
}

func (t *Time) SP99(count int) string {
	return round(t.P99(count))
}

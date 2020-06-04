package metrics

import "time"

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// HistogramData - Histogram에 대한 snapshot 구조 정의
type HistogramData struct {
	Max         int64
	Min         int64
	Mean        float64
	Stddev      float64
	Variance    float64
	Percentiles []float64
}

// Stats - 수집된 Metrics에 대한 snapshot 구조 정의
type Stats struct {
	Time       int64
	Counters   map[string]int64
	Gauges     map[string]int64
	Histograms map[string]HistogramData
}

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewStats - Stats 인스턴스 생성
func NewStats() Stats {
	return Stats{
		Time:       time.Now().UnixNano(),
		Counters:   map[string]int64{},
		Gauges:     map[string]int64{},
		Histograms: map[string]HistogramData{},
	}
}

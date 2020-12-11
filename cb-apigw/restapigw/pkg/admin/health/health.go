// Package health - Health Check 기능 제공 패키지
package health

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
)

// ===== [ Constants and Variables ] =====

const (
	statusOK                 = "OK"
	statusPartiallyAvailable = "Partially Available"
	statusUnavailable        = "Unavailable"
	failureTimeout           = "Timeout during health check"
)

var (
	mu     sync.Mutex
	checks []Config
)

// ===== [ Types ] =====
type (
	// CheckFunc - Check 실행 함수
	CheckFunc func() error

	// Config - Check 처리를 위한 설정 구조
	Config struct {
		Name string
		// 처리 시간 (기본값: 2 second)
		Timeout   time.Duration
		SkipOnErr bool
		Check     CheckFunc
	}

	// Check - Healh Check 상태 정보 구조
	Check struct {
		Status    string            `yaml:"status" json:"status"`
		Timestamp time.Time         `yaml:"timestamp" json:"timestamp"`
		Failures  map[string]string `yaml:"failures,omitempty"`
		System    `yaml:"system" json:"system"`
	}

	// System - Go Process에 대한 정보 구조
	System struct {
		Version          string `yaml:"version" json:"version"`
		GoroutinesCount  int    `yaml:"goroutines_count" json:"goroutines_count"`
		TotalAllocBytes  int    `yaml:"total_alloc_bytes" json:"total_alloc_bytes"`
		HeapObjectsCount int    `yaml:"heap_objects_count" json:"heap_objects_count"`
		AllocBytes       int    `yaml:"alloc_bytes" json:"alloc_bytes"`
	}

	// checkResponse - Check 응답 구조
	checkResponse struct {
		name      string
		skipOnErr bool
		err       error
	}
)

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====

// newSystemMetrics - System 상태 메트릭 정보 구성
func newSystemMetrics() System {
	s := runtime.MemStats{}
	runtime.ReadMemStats(&s)

	return System{
		Version:          runtime.Version(),
		GoroutinesCount:  runtime.NumGoroutine(),
		TotalAllocBytes:  int(s.TotalAlloc),
		HeapObjectsCount: int(s.HeapObjects),
		AllocBytes:       int(s.Alloc),
	}
}

// newCheck - 응답을 기준으로 Check 정보 구성
func newCheck(status string, failures map[string]string) Check {
	return Check{
		Status:    status,
		Timestamp: time.Now(),
		Failures:  failures,
		System:    newSystemMetrics(),
	}
}

// setStatus - 응답 상태 설정
func setStatus(status *string, skipOnErr bool) {
	if skipOnErr && *status != statusUnavailable {
		*status = statusPartiallyAvailable
	} else {
		*status = statusUnavailable
	}
}

// ===== [ Public Functions ] =====

// Register - Check 구성 정보 설정
func Register(c Config) {
	mu.Lock()
	defer mu.Unlock()

	if 0 == c.Timeout {
		c.Timeout = time.Second * 2
	}
	checks = append(checks, c)
}

// Reset - 기존에 설정되어 있던 설정 해체
func Reset() {
	mu.Lock()
	defer mu.Unlock()

	checks = []Config{}
}

// HandlerFunc - Check 처리용 HTTP Handler
func HandlerFunc(rw http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	status := statusOK
	total := len(checks)
	failures := make(map[string]string)
	resChan := make(chan checkResponse, total)

	var wg sync.WaitGroup
	wg.Add(total)

	go func() {
		defer close(resChan)
		wg.Wait()
	}()

	for _, c := range checks {
		go func(c Config) {
			defer wg.Done()
			select {
			case resChan <- checkResponse{c.Name, c.SkipOnErr, c.Check()}:
			default:
			}
		}(c)

	loop:
		for {
			select {
			case <-time.After(c.Timeout):
				failures[c.Name] = failureTimeout
				setStatus(&status, c.SkipOnErr)
				break loop
			case res := <-resChan:
				if nil != res.err {
					failures[res.name] = res.err.Error()
					setStatus(&status, res.skipOnErr)
				}
				break loop
			}
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	c := newCheck(status, failures)
	data, err := core.JSONMarshal(c)
	if nil != err {
		rw.WriteHeader(http.StatusInternalServerError)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	code := http.StatusOK
	if status == statusUnavailable {
		code = http.StatusServiceUnavailable
	}
	rw.WriteHeader(code)
	rw.Write(data)
}

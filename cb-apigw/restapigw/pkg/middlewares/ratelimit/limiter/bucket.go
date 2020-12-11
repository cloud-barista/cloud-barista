// Package limiter - Rate Limit 처리용 TokenBucket 구현 패키지
package limiter

import (
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====

const (
	// rateMargin - Rate 계산할 때 허용 가능한 변동율 (1%)
	rateMargin = 0.01
	// infinityDuration - Unlimited 표현을 위한 Int64 최대 값
	infinityDuration time.Duration = 0x7fffffffffffffff
	// nanosec per second
	nanosec = 1e9
)

var (
	// logger - Logging
	logger = logging.NewLogger()
)

// ===== [ Types ] =====

// IClock - TokenBucket의 시간 처리용 인터페이스
type IClock interface {
	// Now - 현재 시각
	Now() time.Time
	// Sleep - Sleep 처리
	Sleep(d time.Duration)
}

// realClock - 표준 시간 함수를 기준으로 Clock 인터페이스 구현을 위한 구조
type realClock struct{}

/**
 * 용어 정의:
 * - Token : Request 1 개를 소비하는데 사용할 단위
 * - Rate : 초당 허용할 토큰 소비 비율 (일반적으로 Capacity와 동일하게 지정)
 * - Capacity : TokenBucket의 최대 운영 가능한 토큰 수
 * - FillInterval : 초당 Nanoseconds를 기준으로 Capacity당 Rate 비율로 Tick으로 처리할 nanoseconds 계산
 * - Quantum : Tick마다 채워질 토큰 수 (기본 값 1)
 * - Tick : 계산된 FillInterval 단위의 nanoseconds를 1 Tick 단위로 설정
 * - LatestTick : 최종 토큰을 처리한 Tick (Tick은 Bucket의 Start time을 기준으로 산정해서 처리)
 *
 * 계산 알고리즘:
 * - TokenBucket에 대해 Rate(Capacity) 를 기준으로 FillInterval 계산
 * - 토큰의 변동(사용 또는 갱신)이 발생되는 메서드 호출되면 Start time 또는 Latest Tick 기준으로 FillInterval에 대한 Tick * Quantum 만큼의 토큰 추가
 */

// Bucket - Rate Limit 처리를 위한 Token Bucket 구조
type Bucket struct {
	// clock - Bucket 운영에 사용할 Clock (별도 지정이 없으면 시스템 Clock 사용)
	clock IClock
	// startTime - Bucket 생성 시각 (Tick 관리용)
	startTime time.Time
	// capacity -Bucket 최대 Token 용량 (Max Limit 관리용)
	capacity int64
	// quantum - FillInterval에 사용할 Token 수 (기본 값 1)
	quantum int64
	// fillInterval - Token을 추가하기 위한 기간 (Tick 단위)
	fillInterval time.Duration
	// mu - 동시성 관리용
	mu sync.Mutex
	// availableTokens - Latest Tick 기준 현재 사용 가능한 Token 수 (음수인 경우는 Token이 추가될 떄까지 대기 중)
	availableTokens int64
	// latestTick - Token을 처리한 최종 Tick 관리용
	latestTick int64
}

// ===== [ Implementations ] =====

/**
 * RealClock Implements
 */

// Now - 현재 시각 (time.Now)
func (realClock) Now() time.Time {
	return time.Now()
}

// Sleep - Sleep 처리 (time.Sleep)
func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

/**
 * Bucket Implements
 */

// currentTick - TokenBucket의 Start time 기준으로 현재 시각까지의 Tick 갯수 계산
func (tb *Bucket) currentTick(now time.Time) int64 {
	return int64(now.Sub(tb.startTime) / tb.fillInterval)
}

// adjustAvailableTokens - 지정된 Tick 정보를 기준으로 LatestTick 대비로 TokenBucket 내에 사용 가능한 토큰 갯수 조정
func (tb *Bucket) adjustAvailableTokens(tick int64) {
	lastTick := tb.latestTick
	tb.latestTick = tick

	// Capacity를 넘는 경우는 제외
	if tb.availableTokens >= tb.capacity {
		return
	}

	// Tick 수만큼 토큰 추가 (Capacity를 넘을 수 없다)
	tb.availableTokens += (tick - lastTick) * tb.quantum
	if tb.availableTokens > tb.capacity {
		tb.availableTokens = tb.capacity
	}

	logger.Debugf("Adjust Available Token on TokenBucket - [Quantum : %d, fillInterval : %d, lastTick: %d, tick: %d, Added Token: %d (max capacity %d)]", tb.quantum, tb.fillInterval, lastTick, tick, (tick-lastTick)*tb.quantum, tb.capacity)
	return
}

// take - 현재 시각을 기준으로 Blocking 없이 TokenBucket에서 지정한 갯수만큼의 토큰을 삭제(사용)하고, 부족한 경우는 토큰이 추가될 때까지의 대기 시간을 반환 (지정한 최대 대기 시간을 초과하는 경우는 즉시 반환)
func (tb *Bucket) take(now time.Time, count int64, maxWait time.Duration) (time.Duration, bool) {
	if 0 >= count {
		return 0, true
	}

	// 지정한 시각을 기준으로 사용 가능한 토큰 수 조정
	tick := tb.currentTick(now)
	tb.adjustAvailableTokens(tick)

	// 사용 가능 토큰이 존재하는 경우는 대기 없이 즉시 반환
	availableCount := tb.availableTokens - count
	if 0 <= availableCount {
		tb.availableTokens = availableCount
		return 0, true
	}

	// 사용 가능 토큰이 부족한 경우는 대기 시간을 계산, 최대 대기 시간을 넘어가면 사용 불가 상태로 즉시 반환
	endTick := tick + (-availableCount+tb.quantum-1)/tb.quantum
	endTime := tb.startTime.Add(time.Duration(endTick) * tb.fillInterval)
	waitTime := endTime.Sub(now)
	if waitTime > maxWait {
		return 0, false
	}

	// 대기 시간과 사용 가능 토큰 수 설정
	tb.availableTokens = availableCount
	return waitTime, true
}

// takeAvailable - 현재 시간 기준으로 지정한 갯수의 토큰을 사용하는 것으로 처리하고 사용 가능한 토큰 갯수 반환
func (tb *Bucket) takeAvailable(now time.Time, count int64) int64 {
	if 0 >= count {
		return 0
	}

	// 현재 시각을 기준으로 사용 가능 토큰 갯수 조정
	tb.adjustAvailableTokens(tb.currentTick(now))
	if 0 >= tb.availableTokens {
		return 0
	}

	// 부족한 경우는 사용할 수 있는 수량으로 처리
	if count > tb.availableTokens {
		count = tb.availableTokens
	}

	// 사용될 갯수만큼 사용 가능 토큰 수 차감
	tb.availableTokens -= count
	return count
}

// available - 현재 시각 기준으로 TokenBucket 내의 사용 가능한 토큰 갯수 반환
func (tb *Bucket) available(now time.Time) int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 현재 시각을 기준으로 사용 가능한 토큰 갯수 조정
	tb.adjustAvailableTokens(tb.currentTick(now))
	return tb.availableTokens
}

// Take - 현재 시각을 기준으로 Blocking 없이 TokenBucket에서 지정한 갯수의 토큰을 가져오고, 부족한 경우는 토큰이 추가될 떄까지의 대기 시간 반환
func (tb *Bucket) Take(count int64) time.Duration {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	d, _ := tb.take(tb.clock.Now(), count, infinityDuration)
	return d
}

// TakeAvailable - 현재 시각을 기준으로 지정한 갯수의 토크을 사용하는 것으로 처리하고 사용 가능한 토큰 갯수 반환
func (tb *Bucket) TakeAvailable(count int64) int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.takeAvailable(tb.clock.Now(), count)
}

// TakeMaxDuration - 현재 시각을 기준으로 지정한 갯수의 토큰을 얻기 위해서 대기해야 하는 시간과 사용 가능한 토큰의 존재 여부 검증 (지정한 최대 대기 시간을 초과하면 처리 불가로 즉시 반환)
func (tb *Bucket) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.take(tb.clock.Now(), count, maxWait)
}

// Wait - TokenBucket 내의 사용 가능한 토큰 갯수를 기준으로 전달된 갯수의 토큰이 사용 가능할 때까지 대기 처리
func (tb *Bucket) Wait(count int64) {
	if d := tb.Take(count); 0 < d {
		tb.clock.Sleep(d)
	}
}

// WaitMaxDuration - TokenBucket 내의 사용 가능한 토큰 갯수를 기준으로 전달된 갯수의 토큰이 사용 가능할 떄까지 대기 처리 (최대 대기 시간이 초과되면 토큰을 삭제하고 삭제될 토큰이 없는 경우는 즉시 반환)
func (tb *Bucket) WaitMaxDuration(count int64, maxWait time.Duration) bool {
	d, ok := tb.TakeMaxDuration(count, maxWait)
	if 0 < d {
		tb.clock.Sleep(d)
	}
	return ok
}

// Rate - TokenBucket에 설정되어 있는 정보를 기준으로 초당 처리될 수 있는 토큰의 소비 가능 비율 반환
func (tb *Bucket) Rate() float64 {
	return nanosec * float64(tb.quantum) / float64(tb.fillInterval)
}

// Available - TokenBucket 내의 사용 가능한 토큰 갯수 반환
func (tb *Bucket) Available() int64 {
	return tb.available(tb.clock.Now())
}

// Capacity - TokenBucket 의 최대 용량 반환
func (tb *Bucket) Capacity() int64 {
	return tb.capacity
}

// ===== [ Private Functions ] =====

// nextQuantum - 전달된 Quantum 보다 큰 최소의 정수 반환
func nextQuantum(q int64) int64 {
	q1 := q * 11 / 10
	if q1 == q {
		q1++
	}
	return q1
}

// ===== [ Public Functions ] =====

// NewBucket - 지정한 최대 용량까지 지정한 기간마다 1개의 토큰을 채우는 TokenBucket 생성 (생성된 TokenBucket은 최대 용량의 토큰이 설정되어 있다)
func NewBucket(fillInterval time.Duration, capacity int64) *Bucket {
	return NewBucketWithClock(fillInterval, capacity, nil)
}

// NewBucketWithRate - 지정한 비율로 최대 용량까지 토큰을 채우는 TokenBucket 생성 (Clock 정확도가 제한되므로 높은 비율을 지정했을 경우 실제 비율은 1% 정도의 차이가 발생할 수 있다)
func NewBucketWithRate(rate float64, capacity int64) *Bucket {
	return NewBucketWithRateAndClock(rate, capacity, nil)
}

// NewBucketWithRateAndClock - NewBucketWithRate와 동일하며, 운영을 위한 Clock 인터페이스 구현 적용된 TokenBucket 생성
func NewBucketWithRateAndClock(rate float64, capacity int64, clock IClock) *Bucket {
	// Tick 단위로 1개의 토큰을 최대 용량까지 처리하는 TokenBucket 생성
	tb := NewBucketWithQuantumAndClock(1, capacity, 1, clock)

	// 지정한 비율로 Quantum과 FillInterval 결정 (1 ~ 양의 정수 기준 최대 값까지 검증)
	for quantum := int64(1); quantum < 1<<50; quantum = nextQuantum(quantum) {
		// 1초에 소비될 수 있는 비율을 기준으로 토큰이 추가되는 기간을 계산한다.
		fillInterval := time.Duration(nanosec * float64(quantum) / rate)
		if 0 >= fillInterval {
			continue
		}
		tb.fillInterval = fillInterval
		tb.quantum = quantum

		// 계산된 Quantum과 FillInterval을 기준으로 허용되는 비율과 지정한 비율의 편차를 검증한다. (허용 편차를 벗어나는 경우는 재 계산)
		if diff := math.Abs(tb.Rate() - rate); diff/rate <= rateMargin {
			logger.Debugf("[Quantum : %d, fillInterval : %d, Specified Rate: %f, Bucket Rate: %f, Capacity: %d]", quantum, fillInterval, rate, tb.Rate(), capacity)
			return tb
		}
	}

	panic("Cannot find suitable quantum for " + strconv.FormatFloat(rate, 'g', -1, 64))
}

// NewBucketWithClock - NewBucket과 동일하며, 운영을 위한 Clock 인터페이스 구현 적용된 TokenBucket 생성
func NewBucketWithClock(fillInterval time.Duration, capacity int64, clock IClock) *Bucket {
	return NewBucketWithQuantumAndClock(fillInterval, capacity, 1, clock)
}

// NewBucketWithQuantum - NewBucket과 유사하지만, 지정한 기간마다 채워지는 토큰의 갯수를 지정한 TokenBucket 생성
func NewBucketWithQuantum(fillInterval time.Duration, capacity, quantum int64) *Bucket {
	return NewBucketWithQuantumAndClock(fillInterval, capacity, quantum, nil)
}

// NewBucketWithQuantumAndClock - NewBucketWithQuantum과 동일하며, 운영을 위한 Clock 인터페이스 구현 적용된 TokenBucket 쌩성 (Clock이 지정되지 않으면 시스템 사용)
func NewBucketWithQuantumAndClock(fillInterval time.Duration, capacity, quantum int64, clock IClock) *Bucket {
	if nil == clock {
		clock = realClock{}
	}
	if 0 >= fillInterval {
		panic("The token bucket fill interval must be positive.")
	}
	if 0 >= capacity {
		panic("The token bucket capacity must be positive.")
	}
	if 0 >= quantum {
		panic("The token bucket quantum must be positive.")
	}

	return &Bucket{
		clock:           clock,
		startTime:       clock.Now(),
		latestTick:      0,
		fillInterval:    fillInterval,
		capacity:        capacity,
		quantum:         quantum,
		availableTokens: capacity,
	}
}

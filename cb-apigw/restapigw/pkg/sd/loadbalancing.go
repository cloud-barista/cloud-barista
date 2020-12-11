// Package sd -
package sd

import (
	"sync/atomic"

	"github.com/cloud-barista/cb-apigw/restapigw/pkg/config"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/core"
	coreRand "github.com/cloud-barista/cb-apigw/restapigw/pkg/core/rand"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/errors"
	"github.com/cloud-barista/cb-apigw/restapigw/pkg/logging"
)

// ===== [ Constants and Variables ] =====
const ()

var (
	// ErrZeroWeight - Weight가 지정되지 않은 경우 오류
	ErrZeroWeight = errors.New("invalid backend, weight 0 given")
	// ErrCannotElectBackend - Weight 적용등의 방법으로 대상 Backend를 선출할 수 없는 경우 오류
	ErrCannotElectBackend = errors.New("cant elect backend")
)

// ===== [ Types ] =====
type (
	// Balancer - Backend Host에 대한 Load Balancing 전략 적용을 위한 인터페이스
	Balancer interface {
		Host() (string, error)
	}

	// balancer - Subscriber와 연계하기 위한 Balancer 구조
	balancer struct {
		subscriber Subscriber
	}

	// roundRobinLB - Roundrobin 정책 구조
	roundRobinLB struct {
		balancer
		counter uint64
	}

	// randomLB - Random 정책 구조
	randomLB struct {
		balancer
		rand func(uint32) uint32
	}

	// weightLB - Weighted 정책 구조
	weightLB struct {
		balancer

		lastIndex     int // 최종 호출 인덱스
		currentWeight int // 현재 호출 가중치
		gcd           int // 서버별로 지정된 가중치들의 최대 공약수
		maxWeight     int // 서버렬로 지정된 가중치 들의 최대 값
	}

	// noBalancer - 처리대상 Load Balancer가 없는 경우 처리용 형식
	noBalancer string
)

// ===== [ Implementations ] =====

// Host - 연계된 Subscriber를 통해 대상 Host 반환
func (rrb *roundRobinLB) Host() (string, error) {
	hosts, err := rrb.hosts()
	if nil != err {
		return "", err
	}
	offset := (atomic.AddUint64(&rrb.counter, 1) - 1) % uint64(len(hosts))
	logging.GetLogger().Debugf("[MIDDLEWARE] Roundrobin LB > Elected host index: %02d, host: %s", offset, hosts[offset].Host)
	return hosts[offset].Host, nil
}

// Host - 연계된 Subscriber를 통해 대상 Host 반환
func (rb *randomLB) Host() (string, error) {
	hosts, err := rb.hosts()
	if nil != err {
		return "", err
	}
	offset := int(rb.rand(uint32(len(hosts))))
	logging.GetLogger().Debugf("[MIDDLEWARE] Random LB >  Elected host index: %02d, host: %s", offset, hosts[offset].Host)
	return hosts[offset].Host, nil
}

// Host - 연계된 Subscriber를 통해 대상 Host 반환
func (wb *weightLB) Host() (string, error) {
	hosts, err := wb.hosts()
	if nil != err {
		return "", err
	}

	if 1 == len(hosts) {
		return hosts[0].Host, nil
	}

	for {
		// 가중치 기준으로 호출 카운트 검증
		wb.lastIndex = (wb.lastIndex + 1) % len(hosts)
		if wb.lastIndex == 0 {
			wb.currentWeight = wb.currentWeight - wb.gcd
			if wb.currentWeight <= 0 {
				wb.currentWeight = wb.maxWeight
				if wb.currentWeight == 0 {
					return "", ErrCannotElectBackend
				}
			}
		}

		// 산출된 인덱스 기준으로 가중치 적용
		weight := hosts[wb.lastIndex].Weight
		if weight >= wb.currentWeight {
			logging.GetLogger().Debugf("[MIDDLEWARE] Weighted LB > Elected host Index: %d, CurrentWeight: %d, GCD: %d, MaxWeight: %d, Host: %s\n", wb.lastIndex, wb.currentWeight, wb.gcd, wb.maxWeight, hosts[wb.lastIndex].Host)
			return hosts[wb.lastIndex].Host, nil
		}
	}
}

// hosts - 관리 중인 Subscriber의 Hosts 반환
func (b *balancer) hosts() ([]*config.HostConfig, error) {
	hosts, err := b.subscriber.Hosts()
	if nil != err {
		return hosts, err
	}
	if 0 >= len(hosts) {
		return hosts, config.ErrNoHosts
	}
	return hosts, nil
}

// Host - Load Balancing 대상이 없는 경우에 처리할 Dummy Host 반환
func (nb noBalancer) Host() (string, error) {
	return string(nb), nil
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewBalancer - 지정한 Subscriber 정보를 기준으로 Load Balancer 생성
func NewBalancer(subscriber Subscriber) Balancer {
	switch subscriber.Mode() {
	case "rr":
		return NewRoundRobinLB(subscriber)
	case "wrr":
		return NewWeightLB(subscriber)
	default:
		return NewRandomLB(subscriber)
	}
}

// NewRoundRobinLB - 지정한 Subscriber 정보를 기준으로 Roundrobin 정책을 적용한 Load Balancer 생성
func NewRoundRobinLB(subscriber Subscriber) Balancer {
	if fc, ok := subscriber.(FixedSubscriber); ok && 1 == len(fc.hosts) {
		return noBalancer(fc.hosts[0].Host)
	}
	return &roundRobinLB{
		balancer: balancer{subscriber: subscriber},
		counter:  0,
	}
}

// NewRandomLB - 난수 생성을 통해서 랜덤으로 처리하는 Load Balancer 생성
func NewRandomLB(subscriber Subscriber) Balancer {
	if fc, ok := subscriber.(FixedSubscriber); ok && 1 == len(fc.hosts) {
		return noBalancer(fc.hosts[0].Host)
	}
	return &randomLB{
		balancer: balancer{subscriber: subscriber},
		rand:     coreRand.Uint32n,
	}
}

// NewWeightLB - 지정한 Subscriber 정보를 기준으로 Weighted 정책을 적용한 Load Balancer 생성
func NewWeightLB(subscriber Subscriber) Balancer {
	if fc, ok := subscriber.(FixedSubscriber); ok && 1 == len(fc.hosts) {
		return noBalancer(fc.hosts[0].Host)
	}

	wb := &weightLB{balancer: balancer{subscriber: subscriber}}

	// 초기화 작업
	wb.lastIndex = -1
	wb.currentWeight = 0

	hosts, err := wb.hosts()
	if nil == err || 0 < len(hosts) {
		gcdNum := hosts[0].Weight
		max := 0
		for _, host := range hosts {
			gcdNum = core.GCD(gcdNum, host.Weight) // 최대 공약수 계산
			if host.Weight >= max {
				max = host.Weight // 최대 가중치 계산
			}
		}

		wb.gcd = gcdNum
		wb.maxWeight = max
	}

	return wb
}

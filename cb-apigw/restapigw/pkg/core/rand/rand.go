// Package rand - Fast Random 기능 제공 패키지
package rand

import (
	"sync"
	"time"
)

// ===== [ Constants and Variables ] =====
const ()

var (
	// rngPool - 랜덤 값 구성을 위한 Pool
	rngPool sync.Pool
)

// ===== [ Types ] =====
type (
	// RNG - 랜덤 값 생성 구조
	RNG struct {
		x uint32
	}
)

// ===== [ Implementations ] =====

// Uint32 - 관리 중인 Random Seed 값을 기준으로 uint32 랜덤 값 생성
func (r *RNG) Uint32() uint32 {
	for 0 == r.x {
		r.x = getRandomUint32()
	}

	x := r.x
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.x = x

	return x
}

// Uint32n - 관리 중인 Random Seed 값을 기준으로 0 부터 지정한 수까지의 uint32 랜덤 값 생성
func (r *RNG) Uint32n(maxN uint32) uint32 {
	x := r.Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

// ===== [ Private Functions ] =====

// getRandomUint32 - 현재의 Nano 시간 기준으로 Uint32 랜덤 값 생성
func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}

// ===== [ Public Functions ] =====

// Uint32 - 램덤 값을 uint32로 반환
func Uint32() uint32 {
	v := rngPool.Get()
	if nil == v {
		v = &RNG{}
	}
	r := v.(*RNG)
	x := r.Uint32()
	rngPool.Put(r)
	return x
}

// Uint32n - 0 부터 지정한 수까지의 랜덤 값을 uint32로 반환
func Uint32n(maxN uint32) uint32 {
	x := Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

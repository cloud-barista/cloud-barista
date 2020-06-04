// Package backend - Rate Limit 정보를 메모리에 관리하는 패키지
package backend

import (
	"context"
	"sync"
	"time"
)

// ===== [ Constants and Variables ] =====
// ===== [ Types ] =====

// MemoryBackend - sync.Map Wrapping을 사용하는 Memory Backend 구조
type MemoryBackend struct {
	data       map[string]interface{}
	lastAccess map[string]time.Time
	mu         *sync.RWMutex
}

// ===== [ Implementations ] =====

// del - 지정한 키들에 대한 정보를 MemoryBackend에서 삭제
func (mb *MemoryBackend) del(key ...string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for _, k := range key {
		delete(mb.data, k)
		delete(mb.lastAccess, k)
	}
}

// manageEvictions - 지정된 TTL 시간이 지난 데이터 추출 및 삭제
func (mb *MemoryBackend) manageEvictions(ctx context.Context, ttl time.Duration) {
	t := time.NewTicker(ttl)

	for {
		keysToDel := []string{}

		select {
		case <-ctx.Done():
			t.Stop()
			return
		case now := <-t.C:
			mb.mu.RLock()
			for k, v := range mb.lastAccess {
				// 현재 기준으로 TTL 이전인 경우는 삭제
				if v.Add(ttl).Before(now) {
					keysToDel = append(keysToDel, k)
				}
			}
			mb.mu.RUnlock()
		}
		mb.del(keysToDel...)
	}
}

// Load - 지정한 키에 해당하는 정보를 MemoryBackend 에서 로드
func (mb *MemoryBackend) Load(key string, f func() interface{}) interface{} {
	// 키에 해당하는 데이터 추출
	mb.mu.RLock()
	v, ok := mb.data[key]
	mb.mu.RUnlock()

	n := now()

	if ok {
		// Access 시간이 없거나 현재 시간 기준으로 Acess 시간이 지난 경우는 Access 시간 재 설정
		go func(t time.Time) {
			mb.mu.Lock()
			if t0, ok := mb.lastAccess[key]; !ok || t.After(t0) {
				mb.lastAccess[key] = t
			}
			mb.mu.Unlock()
		}(n)

		return v
	}

	mb.mu.Lock()
	defer mb.mu.Unlock()

	v, ok = mb.data[key]
	if ok {
		return v
	}

	v = f()
	mb.lastAccess[key] = n
	mb.data[key] = v
	return v
}

// Store - 지정한 키에 해당하는 정보를 MemoryBackend에 저장
func (mb *MemoryBackend) Store(key string, v interface{}) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.lastAccess[key] = now()
	mb.data[key] = v

	return nil
}

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewMemoryBackend - 지정된 시간동안 유지되는 Memory Backend Store 생성
func NewMemoryBackend(ctx context.Context, ttl time.Duration) *MemoryBackend {
	m := &MemoryBackend{
		data:       map[string]interface{}{},
		lastAccess: map[string]time.Time{},
		mu:         new(sync.RWMutex),
	}

	go m.manageEvictions(ctx, ttl)

	return m
}

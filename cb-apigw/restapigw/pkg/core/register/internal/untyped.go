package internal

import "sync"

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Untyped - 형식 지정이 없는 데이터를 관리하는 구조
type Untyped struct {
	data  map[string]interface{}
	mutex *sync.RWMutex
}

// ===== [ Implementations ] =====

// Register - 지정된 키/값을 형식없는 맵으로 등록
func (u *Untyped) Register(name string, v interface{}) {
	u.mutex.Lock()
	u.data[name] = v
	u.mutex.Unlock()
}

// Get - 지정한 이름에 해당하는 형식없는 값 반환
func (u *Untyped) Get(name string) (interface{}, bool) {
	u.mutex.RLock()
	v, ok := u.data[name]
	u.mutex.RUnlock()
	return v, ok
}

// Clone - 관리 중인 형식없는 맵 복제
func (u *Untyped) Clone() map[string]interface{} {
	u.mutex.RLock()
	res := make(map[string]interface{}, len(u.data))
	for k, v := range u.data {
		res[k] = v
	}
	u.mutex.RUnlock()
	return res
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewUntyped - 형식없는 맴 구조 생성
func NewUntyped() *Untyped {
	return &Untyped{
		data:  map[string]interface{}{},
		mutex: &sync.RWMutex{},
	}
}

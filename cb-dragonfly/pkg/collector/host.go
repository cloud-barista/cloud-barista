package collector

import (
	"sync"
)

type HostInfo struct {
	HostMap *map[string]string
	L       *sync.RWMutex
}

func (h HostInfo) GetHostById(hostId string) string {
	h.L.RLock()
	defer h.L.RUnlock()
	return (*h.HostMap)[hostId]
}

func (h HostInfo) AddHost(hostId string) {
	h.L.Lock()
	defer h.L.Unlock()
	(*h.HostMap)[hostId] = hostId
}

func (h HostInfo) DeleteHost(hostArr []string) {
	h.L.Lock()
	defer h.L.Unlock()
	for _, hostId := range hostArr {
		delete(*h.HostMap, hostId)
	}
}

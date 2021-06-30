package metric

import (
	"reflect"

	"github.com/cloud-barista/cb-dragonfly/pkg/util"
)

// Cpu cpu 메트릭
type Cpu struct {
	CpuGuest       float64 `json:"cpu_guest"`
	CpuGuestNice   float64 `json:"cpu_guest_nice"`
	CpuHardIrq     float64 `json:"cpu_hintr"`
	CpuIdle        float64 `json:"cpu_idle"`
	CpuIowait      float64 `json:"cpu_iowait"`
	CpuNice        float64 `json:"cpu_nice"`
	CpuSoftirq     float64 `json:"cpu_sintr"`
	CpuSteal       float64 `json:"cpu_steal"`
	CpuSystem      float64 `json:"cpu_system"`
	CpuUser        float64 `json:"cpu_user"`
	CpuUtilization float64 `json:"cpu_utilization"`
}

type Cpufreq struct {
	CpuSpeed float64 `json:"cpu_speed"`
}

func (c Cpufreq) GetField() []string {
	val := reflect.ValueOf(c)
	return util.GetFields(val)
}
func (c Cpu) GetField() []string {
	val := reflect.ValueOf(c)
	return util.GetFields(val)
}

// memory 메트릭
type Memory struct {
	MemBuffers     float64 `json:"mem_buffers"`
	MemCached      float64 `json:"mem_cached"`
	MemFree        float64 `json:"mem_free"`
	MemShared      float64 `json:"mem_shared"`
	MemTotal       float64 `json:"mem_total"`
	MemUsed        float64 `json:"mem_used"`
	MemUtilization float64 `json:"mem_utilization"`
}

func (m Memory) GetField() []string {
	val := reflect.ValueOf(m)
	return util.GetFields(val)
}

// disk 메트릭
type Disk struct {
	DiskFree        string `json:"disk_free"`
	DiskTotal       string `json:"disk_total"`
	DiskUsed        string `json:"disk_used"`
	DiskUtilization string `json:"disk_utilization"`
	DiskReadBytes   string `json:"kb_read"`
	DiskWriteBytes  string `json:"kb_written"`
	DIskReadIOPS    int64  `json:"ops_read"`
	DIskWriteIOPS   int64  `json:"ops_write"`
}

func (d Disk) GetField() []string {
	val := reflect.ValueOf(d)
	return util.GetFields(val)
}

// diskio 메트릭
type DiskIO struct {
	DiskReadBytes  string `json:"kb_read"`
	DiskWriteBytes string `json:"kb_written"`
	DIskReadIOPS   int64  `json:"ops_read"`
	DIskWriteIOPS  int64  `json:"ops_write"`
}

func (dio DiskIO) GetField() []string {
	val := reflect.ValueOf(dio)
	return util.GetFields(val)
}

// network 메트릭
type Network struct {
	NetBytesIn   int64 `json:"bytes_in"`
	NetBytesOut  int64 `json:"bytes_out"`
	NetPacketIn  int64 `json:"pkts_in"`
	NetPacketOut int64 `json:"pkts_out"`
}

func (n Network) GetField() []string {
	val := reflect.ValueOf(n)
	return util.GetFields(val)
}

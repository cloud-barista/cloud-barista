package metricstore

import (
	"github.com/cloud-barista/cb-dragonfly/pkg/util"
	"reflect"
)

type MetricType string

const (
	CPU     MetricType = "cpu"
	MEMORY  MetricType = "mem"
	SWAP    MetricType = "swap"
	DISK    MetricType = "disk"
	NETWORK MetricType = "net"
)

// cpu 메트릭
type Cpu struct {
	CpuUtilization float64 `json:"cpu_utilization"`
	CpuSystem      float64 `json:"cpu_system"`
	CpuIdle        float64 `json:"cpu_idle"`
	CpuIowait      float64 `json:"cpu_iowait"`
	CpuHardIrq     float64 `json:"cpu_hintr"`
	CpuSoftirq     float64 `json:"cpu_sintr"`
}

func (c Cpu) GetField() []string {
	val := reflect.ValueOf(c)
	return util.GetFields(val)
}

// memory 메트릭
type Memory struct {
	MemUtilization float64 `json:"mem_utilization"`
	MemTotal       float64 `json:"mem_total"`
	MemUsed        float64 `json:"mem_used"`
	MemFree        float64 `json:"mem_free"`
	MemShared      float64 `json:"mem_shared"`
	MemBuffers     float64 `json:"mem_buffers"`
	MemCached      float64 `json:"mem_cached"`
}

func (m Memory) GetField() []string {
	val := reflect.ValueOf(m)
	return util.GetFields(val)
}

// disk 메트릭
type Disk struct {
	DiskUtilization string `json:"disk_utilization"`
	DiskTotal       string `json:"disk_total"`
	DiskUsed        string `json:"disk_used"`
	DiskFree        string `json:"disk_free"`
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

package influxdb

import (
	"github.com/influxdata/influxdb/client/v2"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Buffer - InfluxDB Client에서 사용할 Buffer 구조 정의
type Buffer struct {
	data []client.BatchPoints
	size int
}

// ===== [ Implementations ] =====

// Add - 지정된 BatchPoints 들을 관리하는 데이터로 추가
func (b *Buffer) Add(ps ...client.BatchPoints) {
	b.data = append(b.data, ps...)
	if len(b.data) > b.size {
		b.data = b.data[len(b.data)-b.size:]
	}
}

// Elements - 관리 중인 BatchPoints 들을 반환하고 초기화
func (b *Buffer) Elements() []client.BatchPoints {
	var res []client.BatchPoints
	res, b.data = b.data, []client.BatchPoints{}
	return res
}

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====

// NewBuffer - 지정한 크기의 Buffer 생성
func NewBuffer(size int) *Buffer {
	return &Buffer{
		data: []client.BatchPoints{},
		size: size,
	}
}

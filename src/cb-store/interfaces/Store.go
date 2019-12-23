// CB-Store is a common repository for managing Meta Info of Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
//
// by powerkim@etri.re.kr, 2019.08.


package interfaces

type KeyValue struct {
	Key string
	Value string
}

type Store interface {
	InitDB() (error)
	InitData() (error)
	Put(key string, value string) (error)
	Get(key string) (*KeyValue, error)
	GetList(key string, sortAscend bool) ([]*KeyValue, error)
	Delete(key string) (error)
}

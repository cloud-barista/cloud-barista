package tumblebug

import (
	"fmt"
	"net/http"
)

/* instance of a Connection Info. */
func NewConnection(name string) *Connection {
	return &Connection{
		Model:      Model{Name: name},
		ConfigName: name,
	}
}

/* instance of a Region */
func NewRegion(name string) *Region {
	return &Region{
		Model:      Model{Name: name},
		RegionName: name,
	}
}

// get a connection info.
func (self *Connection) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/connConfig/%s", self.ConfigName), nil, &self)

}

// get a region
func (self *Region) GET() (bool, error) {

	return self.execute(http.MethodGet, fmt.Sprintf("/region/%s", self.RegionName), nil, &self)

}

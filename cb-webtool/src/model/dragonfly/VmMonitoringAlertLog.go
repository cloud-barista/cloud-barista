package dragonfly

import (
	"time"
	// "fmt"
)

//
type VmMonitoringAlertLog struct {
	Id      string    `json:"id"`
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

// type JSONTime struct {
// 	time.Time
// }

// func (t JSONTime)MarshalJSON() ([]byte, error) {
//     //do your serializing here
//     stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
//     return []byte(stamp), nil
// }

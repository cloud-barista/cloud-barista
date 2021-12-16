package mcis

type StatusCountInfo struct {
	CountCreating    int `json:"countCreating"`
	CountFailed      int `json:"countFailed"`
	CountRebooting   int `json:"countRebooting"`
	CountResuming    int `json:"countResuming"`
	CountRunning     int `json:"countRunning"`
	CountSuspended   int `json:"countSuspended"`
	CountSuspending  int `json:"countSuspending"`
	CountTerminated  int `json:"countTerminated"`
	CountTerminating int `json:"countTerminating"`
	CountTotal       int `json:"countTotal"`
	CountUndefined   int `json:"countUndefined"`
}

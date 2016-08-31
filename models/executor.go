package models

const (
	// StateReady  state ready
	StateReady = "Ready"
	// StateRunning  state running
	StateRunning = "Running"
	// StateDeleted  state deleted
	StateDeleted = "Deleted"

	// ResultSuccess  result success
	ResultSuccess = "OK"
	// ResultFailed  result failed
	ResultFailed = "Error"
	// ResultTimeout  result timeout
	ResultTimeout = "Timeout"
)

// ExecutedResult executor operating result
type ExecutedResult struct {
	//	Name   string
	//	Status string
	//	Detail string
	SID       uint64 `json:"stageVersionID"`
	Namespace string `json:"-"`
	Name      string `json:"stageName"`
	Status    string `json:"runResult"`
	Detail    string `json:"detail"`
}

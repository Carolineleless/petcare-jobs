package job

import "time"

type Status string

const (
	StatusPending    Status = "Pending"
	StatusProcessing Status = "Processing"
	StatusDone       Status = "Done"
	StatusFailed     Status = "Failed"
)

type Job struct {
	ID        string
	Type      string
	Payload   []byte
	Status    Status
	Attempts  int
	CreatedAt time.Time
}

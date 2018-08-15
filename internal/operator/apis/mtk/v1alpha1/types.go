package v1alpha1

import "time"

type Phase string

const (
	PhaseRunning  Phase = "RUNNING"
	PhaseFailed   Phase = "FAILED"
	PhaseComplete Phase = "COMPLETE"
)

// AcquiaStatus to report on the Scheduleds status.
type AcquiaStatus struct {
	Phase   Phase     `json:"phase,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	Message string    `json:"message,omitempty"`
}

type AcquiaDatabase struct {
	Site        string `json:"site"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
}

type Docker struct {
	Image string `json:"image"`
}

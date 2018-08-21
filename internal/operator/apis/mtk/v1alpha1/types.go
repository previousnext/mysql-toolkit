package v1alpha1

import "time"

// Phase of the snapshot.
type Phase string

const (
	// PhaseRunning for when the snapshot is running.
	PhaseRunning Phase = "RUNNING"
	// PhaseFailed for when the snapshot has failed.
	PhaseFailed Phase = "FAILED"
	// PhaseComplete for when the snapshot has completed.
	PhaseComplete Phase = "COMPLETE"
)

// AcquiaStatus to report on the Scheduleds status.
type AcquiaStatus struct {
	Phase   Phase     `json:"phase,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	Message string    `json:"message,omitempty"`
}

// AcquiaDatabase a developer wishes to snapshot.
type AcquiaDatabase struct {
	Site        string `json:"site"`
	Environment string `json:"environment"`
	Name        string `json:"name"`
}

// Docker image being built.
type Docker struct {
	Image string `json:"image"`
}

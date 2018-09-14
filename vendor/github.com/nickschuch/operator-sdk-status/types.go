package exec

// Status provides an update on progression of a stage.
type Status string

const (
	// StatusRunning marks the stage as executing an action.
	StatusRunning Status = "RUNNING"
	// StatusFinished marks the stage as a success.
	StatusFinished Status = "FINISHED"
	// StatusFailed marks the stage as failed, allowing for a retry.
	StatusFailed Status = "FAILED"
)

// Step being executed with other steps in sucession.
type StepStatus struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Message string `json:"message"`
}

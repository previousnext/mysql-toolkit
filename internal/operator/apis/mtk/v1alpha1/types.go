package v1alpha1

import sdkstatus "github.com/nickschuch/operator-sdk-status"

// AcquiaStatus to report on the Scheduleds status.
type AcquiaStatus struct {
	Steps []sdkstatus.StepStatus `json:"steps,omitempty"`
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

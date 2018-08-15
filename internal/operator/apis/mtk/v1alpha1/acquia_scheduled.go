package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/dumper"
)

// AcquiaSnapshotScheduledList of AcquiaSnapshotScheduled objects.
type AcquiaSnapshotScheduledList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AcquiaSnapshotScheduled `json:"items"`
}

// AcquiaSnapshotScheduled declares a one off Acquia database snapshot.
type AcquiaSnapshotScheduled struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AcquiaSnapshotScheduledSpec `json:"spec"`
	Status            AcquiaStatus                `json:"status,omitempty"`
}

// AcquiaSnapshotScheduledSpec declares how to perform a one off Acquia database snapshot.
type AcquiaSnapshotScheduledSpec struct {
	Schedule    string         `json:"schedule"`
	Database    AcquiaDatabase `json:"database"`
	Docker      Docker         `json:"docker"`
	Credentials string         `json:"credentials"`
	Config      dumper.Config  `json:"config"`
}

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/previousnext/mysql-toolkit/internal/dumper"
)

// AcquiaSnapshotList of AcquiaSnapshot objects.
type AcquiaSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []AcquiaSnapshot `json:"items"`
}

// AcquiaSnapshot declares a one off Acquia database snapshot.
type AcquiaSnapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              AcquiaSnapshotSpec `json:"spec"`
	Status            AcquiaStatus       `json:"status,omitempty"`
}

// AcquiaSnapshotSpec declares how to perform a one off Acquia database snapshot.
type AcquiaSnapshotSpec struct {
	Database    AcquiaDatabase `json:"database"`
	Docker      Docker         `json:"docker"`
	Credentials string         `json:"credentials"`
	Config      dumper.Config  `json:"config"`
}

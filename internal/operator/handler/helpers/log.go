package helpers

import (
	"github.com/prometheus/common/log"
	uuid "github.com/satori/go.uuid"
)

// NewLogger which reduces vertical space taken by log message code.
func NewLogger(namespace, name string) log.Logger {
	id := uuid.NewV4()
	return log.With("namespace", namespace).With("name", name).With("id", id)
}

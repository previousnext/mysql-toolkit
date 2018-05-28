package backups

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLatest(t *testing.T) {
	// The second item in this slice should be returned as the latest.
	Backups := []Backup{
		{
			ID:        1,
			Completed: 1527488211, // 28 May 2018 06:16:51
			Link:      "http://foo.bar/baz/1",
		},
		{
			ID:        2,
			Completed: 1527488212, // 29 May 2018 06:16:51
			Link:      "http://foo.bar/baz/2",
		},
		{
			ID:        3,
			Completed: 1524982611, // 29 April 2018 06:16:51
			Link:      "http://foo.bar/baz/3",
		},
	}

	actual := findLatest(Backups)
	assert.Equal(t, int64(2), actual.ID, "latest backup returned")
}
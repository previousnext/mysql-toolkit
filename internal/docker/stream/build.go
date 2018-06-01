package stream

import (
	"encoding/json"
	"fmt"
	"io"
)

// BuildStream returned by the Docker daemon.
type BuildStream struct {
	Stream string `json:"stream"`
}

// Build output which was provided by the Docker daemon.
func Build(w io.Writer, stream io.ReadCloser) error {
	d := json.NewDecoder(stream)

	var s *BuildStream

	for {
		if err := d.Decode(&s); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		fmt.Printf(s.Stream)
	}

	return nil
}

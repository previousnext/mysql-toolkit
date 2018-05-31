package codebuilder

import (
	"os"

	"fmt"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

func compile(name string, files ...string) (string, error) {
	file := fmt.Sprintf("%s/%s.zip", os.TempDir(), name)

	// Get a random uuid for this file.
	err := archiver.Zip.Make(file, files)
	if err != nil {
		return file, errors.Wrap(err, "failed to build source")
	}

	return file, nil
}

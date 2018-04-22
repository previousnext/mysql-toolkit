package codebuilder

import (
	"os"

	"fmt"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

func compile(dockerfile, buildspec, database, name string) (string, error) {
	file := fmt.Sprintf("%s/%s.zip", os.TempDir(), name)

	// Get a random uuid for this file.
	err := archiver.Zip.Make(file, []string{
		dockerfile,
		buildspec,
		database,
	})
	if err != nil {
		return file, errors.Wrap(err, "failed to build source")
	}

	return file, nil
}

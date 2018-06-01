package packager

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

const (
	// Dockerfile path inside tar package.
	Dockerfile = "Dockerfile"
	// Database path inside tar package.
	// @todo, Replace with an ARG.
	Database = "db.sql"
)

// TarParams used for the tar package.
type TarParams struct {
	Files []File
}

// File used for the tar package.
type File struct {
	Local  string
	Remote string
}

// Tar files into a package.
func Tar(params TarParams) (string, error) {
	path := fmt.Sprintf("%s/%s.tar", os.TempDir(), uuid.New().String())

	out, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer out.Close()

	tarWriter := tar.NewWriter(out)
	defer tarWriter.Close()

	for _, f := range params.Files {
		file, err := os.Open(f.Local)
		if err != nil {
			return path, err
		}

		stat, err := file.Stat()
		if err != nil {
			return path, err
		}

		header := &tar.Header{
			Name:    f.Remote,
			Size:    stat.Size(),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return path, err
		}

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return path, err
		}

		err = file.Close()
		if err != nil {
			return path, err
		}
	}

	return path, nil
}

package backups

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// Download the latest database backup.
func (c Client) Download(w io.Writer) error {
	backup, err := c.Latest()
	if err != nil {
		return errors.Wrap(err, "failed to lookup latest backup")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sites/%s/envs/%s/dbs/%s/backups/%s/download.json", baseURL, c.site, c.env, c.db, strconv.FormatInt(backup.ID, 10)), nil)
	if err != nil {
		return errors.Wrap(err, "failed to setup http request")
	}

	req.SetBasicAuth(c.user, c.pass)

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to setup http client")
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to uncompress backup")
	}
	defer reader.Close()

	_, err = io.Copy(w, reader)
	if err != nil {
		return errors.Wrap(err, "failed to save backup")
	}

	return nil
}

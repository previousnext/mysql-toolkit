package backups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// List all backups.
func (c Client) List() ([]Backup, error) {
	var backups []Backup

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sites/%s/envs/%s/dbs/%s/backups.json", baseURL, c.site, c.env, c.db), nil)
	if err != nil {
		return backups, errors.Wrap(err, "failed to setup http client")
	}

	req.SetBasicAuth(c.user, c.pass)

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return backups, errors.Wrap(err, "failed to ")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &backups)
	if err != nil {
		return backups, errors.Wrap(err, "failed to unmarshal")
	}

	return backups, nil
}

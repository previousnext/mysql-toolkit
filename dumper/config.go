package dumper

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const DefaultPlaceholder = "SANITIZED"

type Config struct {
	Sanitize Sanitize `yaml:"sanitize" json:"sanitize"`
}

type Sanitize struct {
	Tables []Table `yaml:"tables"      json:"tables"`
}

type Table struct {
	Name   string  `yaml:"name"   json:"name"`
	Fields []Field `yaml:"fields" json:"fields"`
}

type Field struct {
	Name  string `yaml:"name"  json:"name"`
	Value string `yaml:"value" json:"value"`
}

func Load(path string) (Config, error) {
	var cfg Config

	// We don't want to fail if the config file does not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

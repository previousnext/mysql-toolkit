package dumper

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// DefaultPlaceholder used for setting sanitized fields.
const DefaultPlaceholder = "SANITIZED"

// Config for dumping a MySQL database.
type Config struct {
	Sanitize Sanitize `yaml:"sanitize" json:"sanitize"`
}

// Sanitize rules for while dumping a database.
type Sanitize struct {
	Tables []Table `yaml:"tables"      json:"tables"`
}

// Table rules for while dumping a database.
type Table struct {
	Name   string  `yaml:"name"   json:"name"`
	Fields []Field `yaml:"fields" json:"fields"`
}

// Field rules for while dumping a database.
type Field struct {
	Name  string `yaml:"name"  json:"name"`
	Value string `yaml:"value" json:"value"`
}

// Load a config file.
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

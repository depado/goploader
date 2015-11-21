package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration is the structure that contains a configuration.
type Configuration struct {
	Host string `yaml:"host"`
	TLS  bool   `yaml:"tls"`
}

// C is the exported global configuration variable
var C Configuration

// Load loads the given fp (file path) to the C global configuration variable.
func Load(fp string) error {
	var err error
	conf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(conf, &C)
	return err
}

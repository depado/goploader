package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf is the struct containing the configuration of the server
type Conf struct {
	NameServer   string `yaml:"name_server"`
	UploadDir    string `yaml:"upload_dir"`
	DB           string `yaml:"db"`
	Port         int    `yaml:"port"`
	UniURILength int    `yaml:"uniuri_length"`
	LimitSize    int64  `yaml:"limit_size"`
}

// C is the exported global configuration variable
var C Conf

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

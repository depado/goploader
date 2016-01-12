package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Conf is the struct containing the configuration of the server
type Conf struct {
	NameServer   string `yaml:"name_server" form:"name_server" binding:"required"`
	UploadDir    string `yaml:"upload_dir" form:"upload_dir" binding:"required"`
	DB           string `yaml:"db" form:"db" binding:"required"`
	Port         int    `yaml:"port" form:"port" binding:"required"`
	UniURILength int    `yaml:"uniuri_length" form:"uniuri_length" binding:"required"`
	LimitSize    int64  `yaml:"limit_size" form:"limit_size" binding:"required"`
	FullDoc      bool   `yaml:"fulldoc" form:"fulldoc"`
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

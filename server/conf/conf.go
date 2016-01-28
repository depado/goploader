package conf

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Conf is the struct containing the configuration of the server
type Conf struct {
	NameServer   string
	UploadDir    string
	DB           string
	Port         int
	UniURILength int
	SizeLimit    int64
	TimeLimit    time.Duration
	FullDoc      bool
	Debug        bool
	// TotalSizeLimit int64
}

// UnparsedConf contains an unparsed configuration
type UnparsedConf struct {
	NameServer   string `yaml:"name_server" form:"name_server"`
	UploadDir    string `yaml:"upload_dir" form:"upload_dir"`
	DB           string `yaml:"db" form:"db"`
	Port         int    `yaml:"port" form:"port"`
	UniURILength int    `yaml:"uniuri_length" form:"uniuri_length"`
	SizeLimit    int64  `yaml:"size_limit" form:"size_limit"`
	TimeLimit    string `yaml:"time_limit" form:"time_limit"`
	FullDoc      bool   `yaml:"fulldoc" form:"fulldoc"`
	Debug        bool   `yaml:"debug" form:"debug"`
	// TotalSizeLimit int64  `yaml:"total_size_limit" form:"total_size_limit"`
}

// Validate validates that an unparsed configuration is valid.
func (c *UnparsedConf) Validate() map[string]string {
	errors := make(map[string]string)
	if c.NameServer == "" {
		errors["name_server"] = "This field is required."
	}
	if c.UploadDir == "" {
		c.UploadDir = "up/"
	}
	if c.DB == "" {
		c.DB = "goploader.db"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.UniURILength == 0 {
		c.UniURILength = 10
	}
	if c.SizeLimit == 0 {
		c.SizeLimit = 20
	}
	if c.TimeLimit == "" {
		c.TimeLimit = "2h"
	}
	_, err := time.ParseDuration(c.TimeLimit)
	if err != nil {
		errors["time_limit"] = "Bad duration formatting."
	}
	return errors
}

// C is the exported global configuration variable
var C Conf

// Load loads the given fp (file path) to the C global configuration variable.
func Load(fp string) error {
	var err error
	var uc UnparsedConf
	conf, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(conf, &uc); err != nil {
		return err
	}
	tl, err := time.ParseDuration(uc.TimeLimit)
	if err != nil {
		return err
	}
	C = Conf{
		NameServer:   uc.NameServer,
		UploadDir:    uc.UploadDir,
		DB:           uc.DB,
		Port:         uc.Port,
		UniURILength: uc.UniURILength,
		SizeLimit:    uc.SizeLimit,
		TimeLimit:    tl,
		FullDoc:      uc.FullDoc,
		Debug:        uc.Debug,
	}
	return err
}

package conf

import (
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// C is the exported global configuration variable
var C Conf

// Conf is the struct containing the configuration of the server
type Conf struct {
	NameServer    string `yaml:"name_server" form:"name_server"`
	UploadDir     string `yaml:"upload_dir" form:"upload_dir"`
	DB            string `yaml:"db" form:"db"`
	Port          int    `yaml:"port" form:"port"`
	UniURILength  int    `yaml:"uniuri_length" form:"uniuri_length"`
	KeyLength     int    `yaml:"key_length" form:"key_length"`
	SizeLimit     int64  `yaml:"size_limit" form:"size_limit"`
	NoWeb         bool   `yaml:"no_web" form:"no_web"`
	FullDoc       bool   `yaml:"fulldoc" form:"fulldoc"`
	LogLevel      string `yaml:"loglevel" form:"loglevel"`
	Stats         bool   `yaml:"stats" form:"stats"`
	SensitiveMode bool   `yaml:"sensitive_mode" form:"sensitive_mode"`
}

// NewDefault returns a Conf instance filled with default values
func NewDefault() Conf {
	return Conf{
		UploadDir:    "up/",
		DB:           "goploader.db",
		Port:         8080,
		UniURILength: 10,
		SizeLimit:    20,
		KeyLength:    16,
		LogLevel:     "info",
	}
}

// Validate validates that an unparsed configuration is valid.
func (c *Conf) Validate() map[string]string {
	errors := make(map[string]string)
	if c.NameServer == "" {
		errors["name_server"] = "This field is required."
	}
	return errors
}

// FillDefaults fills the zero value fields in the UnparsedConf with default values
func (c *Conf) FillDefaults() error {
	return mergo.Merge(c, NewDefault())
}

// Load loads the given fp (file path) to the C global configuration variable.
func Load(fp string, verbose bool) error {
	var err error
	var conf []byte
	if conf, err = ioutil.ReadFile(fp); err != nil {
		return err
	}
	if err = yaml.Unmarshal(conf, &C); err != nil {
		return err
	}
	if err = C.FillDefaults(); err != nil {
		return err
	}
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		if err = os.Mkdir(fp, 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

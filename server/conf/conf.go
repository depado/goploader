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
	NameServer string `yaml:"name_server" form:"name_server"`
	Port       int    `yaml:"port" form:"port"`
	AppendPort bool   `yaml:"append_port" form:"append_port"`

	ServeHTTPS bool   `yaml:"serve_https" form:"serve_https"`
	SSLCert    string `yaml:"ssl_cert" form:"ssl_cert"`
	SSLPrivKey string `yaml:"ssl_private_key" form:"ssl_private_key"`

	UploadDir    string  `yaml:"upload_dir" form:"upload_dir"`
	DB           string  `yaml:"db" form:"db"`
	UniURILength int     `yaml:"uniuri_length" form:"uniuri_length"`
	KeyLength    int     `yaml:"key_length" form:"key_length"`
	SizeLimit    int64   `yaml:"size_limit" form:"size_limit"`
	DiskQuota    float64 `yaml:"disk_quota" form:"disk_quota"`
	LogLevel     string  `yaml:"loglevel" form:"loglevel"`

	Stats             bool `yaml:"stats" form:"stats"`
	SensitiveMode     bool `yaml:"sensitive_mode" form:"sensitive_mode"`
	NoWeb             bool `yaml:"no_web" form:"no_web"`
	FullDoc           bool `yaml:"fulldoc" form:"fulldoc"`
	AlwaysDownload    bool `yaml:"always_download" form:"always_download"`
	DisableEncryption bool `yaml:"disable_encryption" form:"disable_encryption"`
}

type UnparsedConf struct {
	NameServer string `yaml:"name_server" form:"name_server"`
	Port       int    `yaml:"port" form:"port"`
	AppendPort bool   `yaml:"append_port" form:"append_port"`

	ServeHTTPS bool   `yaml:"serve_https" form:"serve_https"`
	SSLCert    string `yaml:"ssl_cert" form:"ssl_cert"`
	SSLPrivKey string `yaml:"ssl_private_key" form:"ssl_private_key"`

	UploadDir    string  `yaml:"upload_dir" form:"upload_dir"`
	DB           string  `yaml:"db" form:"db"`
	UniURILength int     `yaml:"uniuri_length" form:"uniuri_length"`
	KeyLength    int     `yaml:"key_length" form:"key_length"`
	SizeLimit    int64   `yaml:"size_limit" form:"size_limit"`
	DiskQuota    float64 `yaml:"disk_quota" form:"disk_quota"`
	LogLevel     string  `yaml:"loglevel" form:"loglevel"`

	Stats             bool `yaml:"stats" form:"stats"`
	SensitiveMode     bool `yaml:"sensitive_mode" form:"sensitive_mode"`
	NoWeb             bool `yaml:"no_web" form:"no_web"`
	FullDoc           bool `yaml:"fulldoc" form:"fulldoc"`
	AlwaysDownload    bool `yaml:"always_download" form:"always_download"`
	DisableEncryption bool `yaml:"disable_encryption" form:"disable_encryption"`
}

// NewDefault returns a Conf instance filled with default values
func NewDefault() Conf {
	return Conf{
		UploadDir:    "up/",
		DB:           "goploader.db",
		Port:         8080,
		UniURILength: 10,
		SizeLimit:    20,
		DiskQuota:    0,
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
	if c.ServeHTTPS {
		if c.SSLCert == "" {
			errors["ssl_cert"] = "This field is required if you serve https."
		}
		if c.SSLPrivKey == "" {
			errors["ssl_private_key"] = "This field is required if you serve https."
		}
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
	if _, err := os.Stat(C.UploadDir); os.IsNotExist(err) {
		if err = os.Mkdir(C.UploadDir, 0777); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

package conf

import (
	"io/ioutil"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"

	"github.com/Depado/goploader/server/logger"
	"github.com/Depado/goploader/server/utils"
)

// C is the exported global configuration variable
var C Conf

// Conf is the struct containing the configuration of the server
type Conf struct {
	NameServer   string `yaml:"name_server" form:"name_server"`
	UploadDir    string `yaml:"upload_dir" form:"upload_dir"`
	DB           string `yaml:"db" form:"db"`
	Port         int    `yaml:"port" form:"port"`
	UniURILength int    `yaml:"uniuri_length" form:"uniuri_length"`
	KeyLength    int    `yaml:"key_length" form:"key_length"`
	SizeLimit    int64  `yaml:"size_limit" form:"size_limit"`
	NoWeb        bool   `yaml:"no_web" form:"no_web"`
	FullDoc      bool   `yaml:"fulldoc" form:"fulldoc"`
	Debug        bool   `yaml:"debug" form:"debug"`
	Stats        bool   `yaml:"stats" form:"stats"`
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
	if verbose {
		logger.Info("server", "Loaded configuration file", fp)
		logger.Info("server", "Name Server", C.NameServer)
		logger.Info("server", "Upload Directory", C.UploadDir)
		logger.Info("server", "Database", C.DB)
		logger.Info("server", "Port", C.Port)
		logger.Info("server", "No Web Interface", C.NoWeb)
		logger.Info("server", "Uni URI Length", C.UniURILength)
		logger.Info("server", "Key Length", C.KeyLength)
		logger.Info("server", "Size Limit", C.SizeLimit, "MB")
		logger.Info("server", "Full Documentation", C.FullDoc)
		logger.Info("server", "Debug", C.Debug)
	}
	return utils.EnsureDir(C.UploadDir)
}

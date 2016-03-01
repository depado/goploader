package conf

import (
	"io/ioutil"
	"log"

	"github.com/imdario/mergo"

	"gopkg.in/yaml.v2"
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
		log.Printf("[INFO][System]\tLoaded configuration file %s :\n", fp)
		log.Printf("[INFO][System]\tName Server : %s\n", C.NameServer)
		log.Printf("[INFO][System]\tUpload Directory : %s\n", C.UploadDir)
		log.Printf("[INFO][System]\tDatabase : %s\n", C.DB)
		log.Printf("[INFO][System]\tPort : %v\n", C.Port)
		log.Printf("[INFO][System]\tNo Web Interface : %v\n", C.NoWeb)
		log.Printf("[INFO][System]\tUni URI Length : %v\n", C.UniURILength)
		log.Printf("[INFO][System]\tKey Length : %v\n", C.KeyLength)
		log.Printf("[INFO][System]\tSize Limit : %v Mo\n", C.SizeLimit)
		log.Printf("[INFO][System]\tFull Documentation : %v\n", C.FullDoc)
		log.Printf("[INFO][System]\tDebug : %v\n", C.Debug)
	}
	return err
}

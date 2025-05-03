package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"

	"gopkg.in/yaml.v2"
)

// Configuration is the struct representing a configuration.
type Configuration struct {
	Service string `yaml:"service"`
	Token   string `yaml:"token"`
}

// C is the exported global configuration variable
var C Configuration

// Load loads the given fp (file path) to the C global configuration variable.
func Load() error {
	var err error
	var hd string
	var conf []byte

	if hd, err = homedir.Dir(); err != nil {
		return err
	}

	cdir := hd + "/.config/"
	cf := cdir + "goploader.conf.yml"

	if _, err = os.Stat(cdir); os.IsNotExist(err) {
		log.Printf("Creating %v directory.\n", cdir)
		if err := os.Mkdir(cdir, 0700); err != nil {
			return fmt.Errorf("creating %v directory: %w", cdir, err)
		}
	} else if err != nil {
		return fmt.Errorf("parsing %v directory: %w", cdir, err)
	}
	if _, err = os.Stat(cf); os.IsNotExist(err) {
		log.Printf("Configuration file %v not found. Writing default configuration.\n", cf)
		C.Service = "http://127.0.0.1:8080"
		if conf, err = yaml.Marshal(C); err != nil {
			return fmt.Errorf("marshalling yaml: %w", err)
		}
		return os.WriteFile(cf, conf, 0644)
	} else if err != nil {
		return err
	}

	if conf, err = os.ReadFile(cf); err != nil {
		return err
	}
	return yaml.Unmarshal(conf, &C)
}

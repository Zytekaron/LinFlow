package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"linflow/src/macro"
	"os"
)

type Config struct {
	LinMod int             `yaml:"linmod"`
	Macros []*macro.Config `yaml:"macros"`
}

func Load(configDirs []string) (cfg *Config, err error) {
	var file *os.File
	for i := 0; i < len(configDirs); i++ {
		file, err = os.Open(configDirs[i])
		if err != nil {
			continue
		}

		err = yaml.NewDecoder(file).Decode(&cfg)
		return
	}

	return nil, errors.New("could not find any config files")
}

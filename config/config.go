package config

import (
	"flag"
	"github.com/mcuadros/go-defaults"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var (
	Default        any
	configFilePath string
)

func loadConfigFile() error {
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Can't read config file: %v\n", err)
		return err
	} else {
		log.Infof("Config file loaded")
		return nil
	}
}

func ParseArgs() {
	flag.StringVar(&configFilePath, "c", "./config.yml", "path to config file location")
	flag.Parse()
	log.Infof("Config file path: %v\n", configFilePath)
}

func LoadConfig[T any]() (*T, error) {
	ParseArgs()
	err := loadConfigFile()
	if err != nil {
		return nil, err
	}

	t := new(T)

	if err = viper.Unmarshal(t); err != nil {
		log.Errorf("Can't unmarshal config file: %v\n", err)
		return nil, err
	} else {
		log.Infof("Config file unmarshalled")
	}

	defaults.SetDefaults(t)
	Default = t

	return t, nil
}

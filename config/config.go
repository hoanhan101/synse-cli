package config

import (
	"fmt"
	//"reflect"
	//"strings"

	log "github.com/Sirupsen/logrus"
	//"github.com/fatih/structs"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type config struct {
	Debug bool
	ActiveHost *hostConfig
	Hosts []*hostConfig
}

type hostConfig struct {
	Name string
	Address string
}

// Config is a new variable containing the config object
var Config config

// ConstructConfig takes in the cli context and builds the current config from
// the cascade of configuration sources. It prioritizes configruation options
// from sources in the following order, with top of the list being highest priority.
//
// 	- Run time CLI flags
// 	- Environment variables
// 	- Configuration files
// 		- .synse.yaml in the local directory
// 		- .synse.yaml in the home (~) directory
//
// All fields in the configuration file are optional.
func ConstructConfig(c *cli.Context) error {
	v := readConfigFromFile()

	err := v.Unmarshal(&Config)
	if err != nil {
		return err
	}

	// Add a host for Synse Server running on localhost
	localHost := hostConfig{
		Name: "local",
		Address: "localhost:5000",

	}
	Config.Hosts = append(Config.Hosts, &localHost)
	if Config.ActiveHost == nil {
		Config.ActiveHost = &localHost
	}

	// FIXME: not sure what this did..
	//s := structs.New(&Config)
	//for _, name := range c.GlobalFlagNames() {
	//	if !c.IsSet(name) {
	//		continue
	//	}
	//
	//	field := s.Field(strings.Replace(strings.Title(name), "-", "", -1))
	//
	//	val := reflect.ValueOf(c.Generic(name)).Elem()
	//
	//	var err error
	//	if val.Kind() == reflect.Bool {
	//		err = field.Set(val.Bool())
	//	} else {
	//		err = field.Set(val.String())
	//	}
	//
	//	if err != nil {
	//		fmt.Printf("%v\n", err)
	//	}
	//}

	log.WithFields(log.Fields{
		"config": fmt.Sprintf("%+v", Config),
	}).Debug("final config")

	return nil
}

// We don't care about being unable to read in the config as it is a non-terminal state.
// Log the issue as debug and move on.
func readConfigFromFile() *viper.Viper {
	v := viper.New()
	v.SetConfigName(".synse")
	v.SetConfigType("yaml")

	v.AddConfigPath(".")      // Try local first
	v.AddConfigPath("$HOME/") // Then try home

	// Defaults
	v.SetDefault("debug", false)
	v.SetDefault("hosts", []hostConfig{})

	v.ReadInConfig()

	log.WithFields(log.Fields{
		"file":     v.ConfigFileUsed(),
		"settings": v.AllSettings(),
	}).Debug("loading config")

	return v
}
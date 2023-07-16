package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	configCmd = func() *cobra.Command {
		configCmd := &cobra.Command{
			Use:   "config",
			Short: "Subcommand to handle config admin functionality of this tool",
			Long:  `Provides capabilities such as initializing an initial configuration as well as parsing`,
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		configCmd.AddCommand(initCmd)
		return configCmd
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the configuration for the tool",
		Long: `There are various fields to be filled up in order to run the configuration.
One can try to initialize the configuration in order to quickly get started with it`,
		Run: func(cmd *cobra.Command, args []string) {
			raw, _ := yaml.Marshal(cfg)
			fmt.Println(string(raw))
		},
	}
)

type serverConfig struct {
	Host        string       `yaml:"host"`
	Port        int          `yaml:"port"`
	IngressPath string       `yaml:"ingressPath"`
	RedirectURI string       `yaml:"redirectUri"`
	Cookie      cookieConfig `yaml:"cookie"`
}

type cookieConfig struct {
	HashKey  string `yaml:"hash_key"`
	BlockKey string `yaml:"block_key"`
}

type mySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type adminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type config struct {
	Server      serverConfig `yaml:"server"`
	Admin       adminConfig  `yaml:"admin"`
	MySQLConfig mySQLConfig  `yaml:"db"`
}

func envVarOrDefault(envVar, defaultVal string) string {
	overrideVal, exists := os.LookupEnv(envVar)
	if exists {
		return overrideVal
	}
	return defaultVal
}

func envVarOrDefaultInt(envVar string, defaultVal int) int {
	overrideVal, exists := os.LookupEnv(envVar)
	if exists {
		num, err := strconv.Atoi(overrideVal)
		if err != nil {
			return defaultVal
		}
		return num
	}
	return defaultVal
}

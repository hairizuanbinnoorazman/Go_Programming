package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func initHashGenerate() string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(rand.Int())))
	return hex.EncodeToString(hasher.Sum(nil))
}

var (
	cfgFile string

	cfg = config{
		Server: serverConfig{
			Host: "0.0.0.0",
			Port: 8080,
			Cookie: cookieConfig{
				HashKey:  initHashGenerate(),
				BlockKey: initHashGenerate(),
			},
		},
		Admin: adminConfig{
			Username: "admin",
			Password: "admin",
		},
		MySQLConfig: mySQLConfig{
			Host:     "localhost",
			Port:     3306,
			DBName:   "basic",
			User:     "user",
			Password: "password",
		},
	}

	serviceName = "backend"
	version     = "v0.1.0"
	rootCmd     = func() *cobra.Command {
		rootCmd := &cobra.Command{
			Use:   "backend",
			Short: "Backend to manage application",
			Long:  ``,
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Help()
			},
		}
		rootCmd.AddCommand(versionCmd)
		rootCmd.AddCommand(configCmd())
		rootCmd.AddCommand(migrateCmd())
		rootCmd.AddCommand(serverCmd())
		return rootCmd
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

func main() {
	rootCmd().Execute()
}

func initConfig() {
	configurationFiles := strings.Split(cfgFile, ",")
	for _, cFile := range configurationFiles {
		var readCfg config
		if strings.Contains(cFile, ".yml") || strings.Contains(cFile, ".yaml") {
			raw, err := os.ReadFile(cFile)
			if err != nil {
				fmt.Println("unable to read config file")
				os.Exit(1)
			}
			err = yaml.Unmarshal(raw, &readCfg)
			if err != nil {
				fmt.Println("unable to process config")
				os.Exit(1)
			}
		}
		mergo.Merge(&cfg, readCfg, mergo.WithOverride)
	}
}

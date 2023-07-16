package main

import (
	"embed"
	"errors"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	//go:embed migrations/*
	migrationFS embed.FS

	migrateCmd = func() *cobra.Command {
		migrateCmd := &cobra.Command{
			Use:   "migrate",
			Short: "Runs database migration (if necessary)",
			Long:  `If one utilizes relational databases such as MySQL - that would require usage data schema migration to happen`,
			Run: func(cmd *cobra.Command, args []string) {
				logger := logrus.New()
				logger.Formatter = stackdriver.NewFormatter(
					stackdriver.WithService(serviceName),
					stackdriver.WithVersion(version),
				)
				logger.Level = logrus.InfoLevel
				logger.Info("Run migration")
				defer logger.Info("Migration completed")

				d, err := iofs.New(migrationFS, "migrations")
				if err != nil {
					log.Fatal(err)
				}

				m, err := migrate.NewWithSourceInstance(
					"iofs", d, fmt.Sprintf("mysql://%s:%s@(%s:%d)/%s", cfg.MySQLConfig.User, cfg.MySQLConfig.Password, cfg.MySQLConfig.Host, cfg.MySQLConfig.Port, cfg.MySQLConfig.DBName))
				defer func() {
					_, err := m.Close()
					if err != nil {
						logger.Errorf("unable to disconnect from database correctly :: %v", err)
					}
				}()
				if err != nil {
					panic(fmt.Sprintf("unable to connect to database :: %v", err))
				}
				err = m.Up()
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					panic(fmt.Sprintf("unable to upgrade schema :: %v", err))
				}
			},
		}
		migrateCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Configuration File")
		return migrateCmd
	}
)

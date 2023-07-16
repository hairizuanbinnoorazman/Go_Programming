package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/hairizuanbinnoorazman/basic-app/cmd/backend/handlers"
)

var (
	serverCmd = func() *cobra.Command {
		serverCmd := &cobra.Command{
			Use:   "server",
			Short: "Run the API server",
			Long:  `Runs the API server`,
			Run: func(cmd *cobra.Command, args []string) {
				logger := logrus.New()
				logger.Formatter = stackdriver.NewFormatter(
					stackdriver.WithService(serviceName),
					stackdriver.WithVersion(version),
				)
				logger.Level = logrus.InfoLevel
				logger.Info("Application Start Up")
				defer logger.Info("Application Ended")

				connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True", cfg.MySQLConfig.User, cfg.MySQLConfig.Password, cfg.MySQLConfig.Host, cfg.MySQLConfig.Port, cfg.MySQLConfig.DBName)
				_, err := gorm.Open(mysql.Open(connectionString))
				if err != nil {
					logger.Errorf("Unable to create mysql client. %v", err)
					os.Exit(1)
				}

				cookieAuth := handlers.Auth{
					HashKey:    []byte(cfg.Server.Cookie.HashKey),
					BlockKey:   []byte(cfg.Server.Cookie.BlockKey),
					CookieName: "basic",
				}

				_ = handlers.AuthWrapper{
					Auth:   cookieAuth,
					Logger: logger,
				}

				r := mux.NewRouter()
				r.NotFoundHandler = NotFound{
					Logger: logger,
				}
				r.Handle("/status", handlers.Status{
					Logger: logger,
				})
				r.Handle("/healthz", handlers.Status{
					Logger: logger,
				})
				r.Handle("/readyz", handlers.Status{
					Logger: logger,
				})

				srv := http.Server{
					Handler:      r,
					Addr:         fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port),
					WriteTimeout: 15 * time.Second,
					ReadTimeout:  15 * time.Second,
				}

				addFrontendRoutes(r)

				logger.Fatal(srv.ListenAndServe())
			},
		}
		serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Configuration File")
		return serverCmd
	}
)

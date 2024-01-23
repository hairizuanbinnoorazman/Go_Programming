package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	gormMySQL "gorm.io/driver/mysql"
	gormSQLite "gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"embed"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
)

//go:embed migrations/*
var fs embed.FS

type User struct {
	ID        int `gorm:"primaryKey,autoIncrement"`
	FirstName string
	LastName  string
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "This is a sample golang migrate application",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migration",
	Run: func(cmd *cobra.Command, args []string) {
		dbUser := os.Getenv("DATABASE_USER")
		dbPass := os.Getenv("DATABASE_PASSWORD")
		dbHost := os.Getenv("DATABASE_HOST")
		dbName := os.Getenv("DATABASE_NAME")
		useTLS := os.Getenv("DATABASE_USE_TLS")
		dbType := os.Getenv("DATABASE_TYPE")

		var d source.Driver
		var err error
		dsn := ""
		if dbType == "" || dbType == "mysql" {
			dsn = fmt.Sprintf("mysql://%v:%v@(%v:3306)/%v", dbUser, dbPass, dbHost, dbName)
			if strings.ToLower(useTLS) == "true" {
				fmt.Println("database tls mode on")
				dsn = dsn + "?tls=true"
			}
			d, err = iofs.New(fs, "migrations/mysql")
			if err != nil {
				log.Fatal(err)
			}
		} else if dbType == "sqlite" {
			sqliteFile := os.Getenv("SQLITE_FILE")
			dsn = fmt.Sprintf("sqlite3://%s?query", sqliteFile)
			d, err = iofs.New(fs, "migrations/sqlite")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("unexpected dbType provided. Please check inputs")
			os.Exit(1)
		}

		m, err := migrate.NewWithSourceInstance(
			"iofs", d, dsn)

		if err != nil {
			panic(fmt.Sprintf("unable to connect to database :: %v", err))
		}
		err = m.Up()
		if err == migrate.ErrNoChange {
			fmt.Println("no change to database")
			os.Exit(0)
		}
		if err != nil {
			panic(fmt.Sprintf("unable to connect to database :: %v", err))
		}
	},
}

type Status struct{}

func (h Status) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start status endpoint")
	defer log.Println("end status endpoint")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type UserGet struct {
	DB *gorm.DB
}

func (h UserGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawUserID := vars["userID"]
	userID, err := strconv.Atoi(rawUserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	var u User
	result := h.DB.First(&u, userID)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("bad connection"))
		return
	}

	rawResp, _ := json.Marshal(u)
	w.WriteHeader(http.StatusOK)
	w.Write(rawResp)
}

type UserCreate struct {
	DB *gorm.DB
}

func (h UserCreate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	type userCreate struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var uc userCreate
	json.Unmarshal(raw, &uc)

	u1 := User{FirstName: uc.FirstName, LastName: uc.LastName}
	result := h.DB.Create(&u1)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("bad connection"))
		return
	}

	rawResp, _ := json.Marshal(u1)
	w.WriteHeader(http.StatusOK)
	w.Write(rawResp)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server start")
		dbUser := os.Getenv("DATABASE_USER")
		dbPass := os.Getenv("DATABASE_PASSWORD")
		dbHost := os.Getenv("DATABASE_HOST")
		dbName := os.Getenv("DATABASE_NAME")
		useTLS := os.Getenv("DATABASE_USE_TLS")
		dbType := os.Getenv("DATABASE_TYPE")

		var dsn string
		var db *gorm.DB
		var err error
		if dbType == "" || dbType == "mysql" {
			dsn = fmt.Sprintf("%v:%v@(%v:3306)/%v", dbUser, dbPass, dbHost, dbName)
			if strings.ToLower(useTLS) == "true" {
				fmt.Println("database tls mode on")
				dsn = dsn + "?tls=true"
			}
			db, err = gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
		} else if dbType == "sqlite" {
			sqliteFile := os.Getenv("SQLITE_FILE")
			dsn = fmt.Sprintf("%s?query", sqliteFile)
			db, err = gorm.Open(gormSQLite.Open(dsn), &gorm.Config{})
		}

		if err != nil {
			panic(fmt.Sprintf("unable to connect to database :: %v", err))
		}

		r := mux.NewRouter()
		r.Handle("/health", Status{}).Methods("GET")
		r.Handle("/user", UserCreate{DB: db}).Methods("POST")
		r.Handle("/user/{userID}", UserGet{DB: db}).Methods("GET")

		srv := &http.Server{
			Handler: r,
			Addr:    ":8888",
		}

		log.Fatal(srv.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

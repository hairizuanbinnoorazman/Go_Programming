package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"embed"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
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
		d, err := iofs.New(fs, "migrations")
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithSourceInstance(
			"iofs", d, "mysql://user:password@(localhost:3306)/application")

		if err != nil {
			panic(fmt.Sprintf("unable to connect to database :: %v", err))
		}
		m.Up()
	},
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
	raw, err := ioutil.ReadAll(r.Body)
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
		dsn := "user:password@tcp(127.0.0.1:3306)/application"
		db, err := gorm.Open(gormMySQL.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("unable to connect to database :: %v", err))
		}

		r := mux.NewRouter()
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

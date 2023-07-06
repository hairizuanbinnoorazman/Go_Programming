package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"
)

var (
	serviceName = "shopping-list"
	version     = "1.0.0"
	entityTable = "ShoppingItems"
)

type Auth struct {
	HashKey    []byte
	BlockKey   []byte
	CookieName string
}

func main() {
	logger := logrus.New()
	logger.Formatter = stackdriver.NewFormatter(
		stackdriver.WithService(serviceName),
		stackdriver.WithVersion(version),
	)
	logger.Level = logrus.InfoLevel
	logger.Info("Server start")

	adminUser := os.Getenv("ADMIN_USER")
	if adminUser == "" {
		log.Println("no admin user set")
		os.Exit(1)
	}

	adminPassword := os.Getenv("ADMIN_PASS")
	if adminPassword == "" {
		log.Println("no admin pass set")
		os.Exit(1)
	}

	hashKey := os.Getenv("HASH_KEY")
	if hashKey == "" {
		log.Println("no admin pass set")
		os.Exit(1)
	}

	blockKey := os.Getenv("BLOCK_KEY")
	if blockKey == "" {
		log.Println("no admin pass set")
		os.Exit(1)
	}

	cookieName := os.Getenv("COOKIE_NAME")
	if cookieName == "" {
		log.Println("no admin pass set")
		os.Exit(1)
	}

	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if cookieName == "" {
		log.Println("no project id set")
		os.Exit(1)
	}

	dsClient, err := datastore.NewClient(context.TODO(), projectID)
	if err != nil {
		logger.Fatalf("unable to connect to datastore")
		panic("failed")
	}

	r := mux.NewRouter()
	r.Handle("/login", Login{Logger: logger, AdminUser: adminUser, AdminPass: adminPassword})
	r.Handle("/api/shopping-list/v1/item", AddShoppingItem{
		Logger:      logger,
		Datastore:   *dsClient,
		EntityTable: entityTable,
	}).Methods("POST")
	r.Handle("/api/shopping-list/v1/item", ListShoppingItems{
		Logger:      logger,
		Datastore:   *dsClient,
		EntityTable: entityTable,
	}).Methods("GET")
	srv := http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%v:%v", "0.0.0.0", "8080"),
	}

	logger.Fatal(srv.ListenAndServe())
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type Login struct {
	Logger    Logger
	AdminUser string
	AdminPass string
	Auth      Auth
}

func (h Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("start Login handler")
	defer h.Logger.Infof("end Login handler")

	type loginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("died"))
		return
	}
	var l loginInput
	json.Unmarshal(raw, &l)
	if l.Username != h.AdminUser || l.Password != h.AdminPass {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("bad username or password"))
		return
	}

	value := map[string]string{
		"username": h.AdminUser,
	}

	s := securecookie.New(h.Auth.HashKey, h.Auth.BlockKey)
	encoded, err := s.Encode(h.Auth.CookieName, value)
	if err != nil {
		errMsg := fmt.Sprintf("Error - unable to set authorization token. Error: %v", err)
		h.Logger.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}
	cookie := &http.Cookie{
		Name:     h.Auth.CookieName,
		Value:    encoded,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type AuthDecorator struct {
	Logger      Logger
	Auth        Auth
	AdminUser   string
	NextHandler http.Handler
}

func (h AuthDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("start AuthDecorator handler")
	defer h.Logger.Infof("end AuthDecorator handler")

	ctx := r.Context()
	s := securecookie.New(h.Auth.HashKey, h.Auth.BlockKey)
	cookie, cookieErr := r.Cookie(h.Auth.CookieName)
	value := make(map[string]string)
	if cookieErr == nil {
		err := s.Decode(h.Auth.CookieName, cookie.Value, &value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			h.Logger.Errorf("unable to decode cookie : %v", err)
			return
		}
	} else {
		h.Logger.Error(cookieErr)
	}

	if value["username"] != h.AdminUser {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		h.Logger.Errorf("username is not the same")
		return
	}

	h.NextHandler.ServeHTTP(w, r.WithContext(ctx))
}

type status string

var (
	active status = "active"
	incart status = "in-cart"
)

type ShoppingItem struct {
	Name    string
	Created time.Time
	Status  status
}

type AddShoppingItem struct {
	Logger      Logger
	Datastore   datastore.Client
	EntityTable string
}

func (h AddShoppingItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("start AddShoppingItem handler")
	defer h.Logger.Infof("end AddShoppingItem handler")
	ctx := r.Context()

	type addShoppingListReq struct {
		Name string `json:"name"`
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorf("unable to read request body :: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}
	var req addShoppingListReq
	json.Unmarshal(raw, &req)

	uid, _ := uuid.NewV4()

	key := datastore.NameKey(h.EntityTable, uid.String(), nil)
	e := ShoppingItem{Name: req.Name, Created: time.Now(), Status: active}
	_, err = h.Datastore.Put(ctx, key, &e)

	if err != nil {
		h.Logger.Errorf("unable to save data into datastore :: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("created"))
}

type ModifyShoppingItem struct {
	Logger    Logger
	Datastore datastore.Client
}

func (h ModifyShoppingItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("start ModifyShoppingItem handler")
	defer h.Logger.Infof("end ModifyShoppingItem handler")
}

type DeleteShoppingItem struct {
	Logger    Logger
	Datastore datastore.Client
}

func (h DeleteShoppingItem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("start DeleteShoppingItem handler")
	defer h.Logger.Infof("end DeleteShoppingItem handler")
}

type ListShoppingItems struct {
	Logger      Logger
	Datastore   datastore.Client
	EntityTable string
}

func (h ListShoppingItems) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("start ListShoppingItems handler")
	defer h.Logger.Infof("end ListShoppingItems handler")

	ctx := r.Context()

	var items []ShoppingItem
	q := datastore.NewQuery(h.EntityTable)
	_, err := h.Datastore.GetAll(ctx, q, &items)

	if err != nil {
		h.Logger.Errorf("unable to get all records from datastore :: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	resp := map[string]interface{}{}
	resp["items"] = items
	respRaw, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(respRaw)
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/hairizuanbinnoorazman/basic-app/logger"
)

type cookieKey string
type authType string

var (
	userIDKey     cookieKey = "userID"
	userTypeKey   cookieKey = "userType"
	AdminUserType authType  = "admin"
	UserType      authType  = "user"
)

type Auth struct {
	// Secret     string
	// ExpiryTime int
	// Issuer     string
	HashKey    []byte
	BlockKey   []byte
	CookieName string
}

type AuthWrapper struct {
	Auth   Auth
	Logger logger.Logger
}

func (w AuthWrapper) RequireAuth(minLevel authType, handler http.Handler) AuthDecorator {
	return AuthDecorator{
		MinAuthLevel: minLevel,
		Auth:         w.Auth,
		Logger:       w.Logger,
		NextHandler:  handler,
	}
}

type AuthDecorator struct {
	MinAuthLevel authType
	Auth         Auth
	Logger       logger.Logger
	NextHandler  http.Handler
}

func (a AuthDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Logger.Info("RequireAuth Exists Check")

	ctx := r.Context()
	s := securecookie.New(a.Auth.HashKey, a.Auth.BlockKey)
	value := make(map[cookieKey]string)
	var userID string
	var userType authType

	type failedResp struct {
		Msg string `json:"msg"`
	}
	rawErrMsg, _ := json.Marshal(failedResp{Msg: "Invalid authorization token"})

	cookie, cookieErr := r.Cookie(a.Auth.CookieName)
	if cookieErr == nil {
		err := s.Decode(a.Auth.CookieName, cookie.Value, &value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(rawErrMsg)
			a.Logger.Errorf("unable to decode cookie : %v", err)
			return
		}
		userID = value[userIDKey]
		userType = authType(value[userTypeKey])
	} else {
		a.Logger.Error(cookieErr)
	}

	if userType == "" || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(rawErrMsg)
		a.Logger.Info("UserID or UserType is empty/bad after checking cookie and header")
		return
	}

	if a.MinAuthLevel == AdminUserType {
		if userType != AdminUserType {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(rawErrMsg)
			a.Logger.Infof("%v attempted to login for %v", userType, r.URL.Path)
			return
		}
	}

	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, userTypeKey, userType)

	a.NextHandler.ServeHTTP(w, r.WithContext(ctx))
}

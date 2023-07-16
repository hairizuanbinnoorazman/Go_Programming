package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/hairizuanbinnoorazman/basic-app/logger"
	"github.com/hairizuanbinnoorazman/basic-app/user"
)

// Login - Handles situation of user signing in if user exists - else, throw back error
// Require search of user by email address
type Login struct {
	Logger      logger.Logger
	UserStore   user.Store
	Auth        Auth
	RedirectURI string
}

func (h Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Start Login Handler")
	defer h.Logger.Info("End Login Handler")

	rawReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Error - unable to read json body. Error: %v", err)
		h.Logger.Error(errMsg)
		w.WriteHeader(400)
		w.Write([]byte(generateErrorResp(errMsg)))
		return
	}

	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := loginReq{}
	err = json.Unmarshal(rawReq, &req)
	if err != nil {
		errMsg := fmt.Sprintf("Error - unable to parse login body. Error: %v", err)
		h.Logger.Error(errMsg)
		w.WriteHeader(400)
		w.Write([]byte(generateErrorResp(errMsg)))
		return
	}

	u, err := h.UserStore.GetUserByEmail(context.TODO(), req.Email)
	if err != nil {
		errMsg := fmt.Sprintf("Error - unable to find user. Error: %v", err)
		h.Logger.Error(errMsg)
		w.WriteHeader(404)
		w.Write([]byte(generateErrorResp(errMsg)))
		return
	}

	passwordCorrect := u.IsPasswordCorrect(req.Password)
	if !passwordCorrect {
		errMsg := fmt.Sprintf("Error - unable to find user. Error: %v", err)
		h.Logger.Error(errMsg)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(generateErrorResp(errMsg)))
		return
	}

	value := map[string]string{
		"user_id": u.ID,
	}
	s := securecookie.New(h.Auth.HashKey, h.Auth.BlockKey)
	encoded, err := s.Encode(h.Auth.CookieName, value)
	if err != nil {
		errMsg := fmt.Sprintf("Error - unable to set authorization token. Error: %v", err)
		h.Logger.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(generateErrorResp(errMsg)))
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
	w.Write([]byte(""))
}

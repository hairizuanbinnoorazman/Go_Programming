package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

// Hash keys should be at least 32 bytes long
var hashKey = []byte("very-secret test test test")

// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
// Shorter keys may weaken the encryption used.
var blockKey = []byte("this needs to be a very long sec")
var s = securecookie.New(hashKey, blockKey)

func main() {
	fmt.Println("Start server")

	http.HandleFunc("/set", SetCookieHandler)
	http.HandleFunc("/read", ReadCookieHandler)

	http.ListenAndServe(":8080", nil)
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	value := map[string]string{
		"foo": "bar",
	}
	encoded, err := s.Encode("cookie-name", value)

	if err == nil {
		cookie := &http.Cookie{
			Name:     "cookie-name",
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	} else {
		fmt.Printf("error trying to encode and set cookie. Err: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is a setter page"))
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("cookie-name")
	if err == nil {
		value := make(map[string]string)
		if err = s.Decode("cookie-name", cookie.Value, &value); err == nil {
			fmt.Fprintf(w, "The value of foo is %q", value["foo"])
		}
	} else {
		fmt.Printf("No cookie found. Err: %v\n", err)
	}

	w.Write([]byte("Check the server logs"))
}

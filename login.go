package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/sha3"
)

func login(hash string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		userpass := base64.URLEncoding.EncodeToString(sha3.New512().Sum([]byte(r.FormValue("password"))))
		secret := r.FormValue("secret")

		if userpass == hash {
			http.SetCookie(w, &http.Cookie{Name: "secret", Value: hash, Expires: time.Now().Add(time.Hour * 24 * 31 * 12), Path: "/"})
			json.NewEncoder(w).Encode(true)
			return
		}
		if secret == hash {
			json.NewEncoder(w).Encode(true)
		} else {
			json.NewEncoder(w).Encode(false)
		}

	}
}

type passwordFlag struct {
	password string
}

func (pass *passwordFlag) String() string {
	return fmt.Sprint(pass.password)
}

func (pass *passwordFlag) Set(value string) error {
	pass.password = base64.URLEncoding.EncodeToString(sha3.New512().Sum([]byte(value)))
	return nil
}

func authorize(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("secret")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}
	if cookie.Value != passArg.password {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("Password and secret do not match")
	}
	return nil
}

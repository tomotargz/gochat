package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf = &oauth2.Config{
	ClientID:     os.Getenv("GOCHAT_CLIENT_ID"),
	ClientSecret: os.Getenv("GOCHAT_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/login/google/callback",
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world!")
	})

	http.HandleFunc("/login/google/oauth2", func(w http.ResponseWriter, r *http.Request) {
		url := conf.AuthCodeURL("state")
		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc("/login/google/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome back!")
		// TODO: print user name
	})

	http.ListenAndServe(":8080", nil)
}

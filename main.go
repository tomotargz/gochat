package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf = oauth2.Config{
	ClientID:     os.Getenv("GOCHAT_CLIENT_ID"),
	ClientSecret: os.Getenv("GOCHAT_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/login/google/callback",
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,
}

type userInfo struct {
	Name string `json:"name"`
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
		ctx := context.Background()
		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		client := conf.Client(ctx, token)
		res, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		byteArray, err := ioutil.ReadAll(res.Body)
		user := userInfo{}
		json.Unmarshal(byteArray, &user)
		fmt.Fprintf(w, "Hello %s!", user.Name)
	})

	http.ListenAndServe(":8080", nil)
}

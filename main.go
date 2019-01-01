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

var oauth2Conf = oauth2.Config{
	ClientID:     os.Getenv("GOCHAT_CLIENT_ID"),
	ClientSecret: os.Getenv("GOCHAT_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,
}

type userInfo struct {
	Name string `json:"name"`
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World.")
}

func login(w http.ResponseWriter, r *http.Request) {
	url := oauth2Conf.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	token, err := oauth2Conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	client := oauth2Conf.Client(ctx, token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	byteArray, err := ioutil.ReadAll(res.Body)
	user := userInfo{}
	json.Unmarshal(byteArray, &user)
	fmt.Fprintf(w, "Hello %s!", user.Name)
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)
	http.ListenAndServe(":8080", nil)
}

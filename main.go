package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var sessions = make(map[string]string)

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

type message struct {
	User    string
	Comment string
}

var timeline []message

type chatTemplateSource struct {
	User     string
	Messages []message
}

func root(w http.ResponseWriter, r *http.Request) {
	name, ok := auth(r)
	if !ok {
		http.Redirect(w, r, oauth2Conf.AuthCodeURL("state"), http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("chat.html")

	source := chatTemplateSource{
		User:     name,
		Messages: timeline,
	}
	t.Execute(w, source)
}

func generateHash(s string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	e := strings.Split(string(b), "$")
	saltAndHash := e[len(e)-1]
	hash := string([]rune(saltAndHash)[22:])
	return hash
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
	session := generateHash(user.Name)
	sessions[session] = user.Name
	c := &http.Cookie{
		Name:  "SESSION",
		Value: session,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/callback", callback)
	room := newRoom()
	go room.run()
	http.Handle("/ws", room)
	http.ListenAndServe(":8080", nil)
}

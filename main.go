package main

import (
	"fmt"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World.")
}

func main() {
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}

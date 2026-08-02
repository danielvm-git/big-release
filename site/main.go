package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
)

//go:embed index.html
var indexHTML []byte

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(indexHTML); err != nil {
		return
	}
}

func listenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(listenAddr(), nil))
}

package main

import (
	"bytes"
	"log"
	"log/slog"
	"net/http"
)

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Welcome to our Homepage!\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
}

func handleGoodbye(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Goodbye, World!\n"))
	if err != nil {
		slog.Error("error writing response", "err", err)
		return
	}
}

func handleHelloParameterized(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userList := params["user"]

	username := "User"
	if len(userList) > 0 {
		username = userList[0]
	}

	var output bytes.Buffer
	output.WriteString("Hello, ")
	output.WriteString(username)
	output.WriteString("!\n")

	_, err := w.Write(output.Bytes())
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/goodbye", handleGoodbye)
	mux.HandleFunc("/hello/", handleHelloParameterized)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

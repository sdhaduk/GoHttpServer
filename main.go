package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type UserData struct {
	Name string
}

func buildOutput(w http.ResponseWriter, username string) error {
	var output bytes.Buffer

	output.WriteString("Hello, ")
	output.WriteString(username)
	output.WriteString("!\n")

	_, err := w.Write(output.Bytes())
	if err != nil {
		return err
	}

	return nil
}

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

	err := buildOutput(w, username)
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func handleUserResponsesHello(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")
	
	err := buildOutput(w, username)
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func handleHelloHeader(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user")

	if username == "" {
		http.Error(w, "invalid username provided", http.StatusBadRequest)
		return
	}
	
	err := buildOutput(w, username)
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	byteData, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading request body", "err", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}
	var reqData UserData
	err = json.Unmarshal(byteData, &reqData)
	if err != nil {
		slog.Error("error unmarshalling request body", "err", err)
		http.Error(w, "error parsing request JSON", http.StatusBadRequest)
		return
	}

	err = buildOutput(w, reqData.Name)
	if err != nil {
		slog.Error("error writing response body", "err", err)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/goodbye", handleGoodbye)
	mux.HandleFunc("/hello/", handleHelloParameterized)
	mux.HandleFunc("/responses/{user}/hello/", handleUserResponsesHello)
	mux.HandleFunc("/user/hello/", handleHelloHeader)
	mux.HandleFunc("/json", handleJSON)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

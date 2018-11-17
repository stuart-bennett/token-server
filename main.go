package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type tokenStore map[string]string

const LoginPath string = "/login"
const UsernamePath string = "/username"

func main() {
	log.Fatal(http.ListenAndServe("localhost:8000", ConfigureMux()))
}

func ConfigureMux() *http.ServeMux {
	ts := tokenStore{}
	m := http.NewServeMux()
	m.Handle(LoginPath, http.HandlerFunc(ts.login))
	m.Handle(UsernamePath, http.HandlerFunc(ts.username))
	return m
}

func (ts tokenStore) login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var ltr LoginTokenRequest;
	if err := json.NewDecoder(req.Body).Decode(&ltr); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if authenticate(ltr.Username, ltr.Password) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{ \"token\": \"temp-token\" }")
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func authenticate(u string, p string) bool {
	return u == "admin" && p == "admin1000"
}

func (ts tokenStore) username(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type LoginTokenRequest struct {
	Username string;
	Password string;
}

type LoginTokenResponse struct {
	Token string
}

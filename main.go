package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type tokenStore map[string]string

const TokenRequestHeader string = "X-Auth-Token"
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

	ltr, ok := validateLoginTokenRequest(req)
	if !ok {
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

func validateLoginTokenRequest(req *http.Request) (LoginTokenRequest, bool) {
	var ltr LoginTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&ltr); err != nil {
		log.Printf("Could not decode request because %s", err)
		return LoginTokenRequest{}, false
	}

	if ltr.Username == "" || ltr.Password == "" {
		return LoginTokenRequest{}, false
	}

	return ltr, true
}

func authenticate(u string, p string) bool {
	return u == "admin" && p == "admin1000"
}

func verifyToken(t string) bool {
	return t == "test-token"
}

func (ts tokenStore) username(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	token, ok := validateUsernameRequest(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if verifyToken(token) {
		resp, err := json.Marshal(UsernameResponse{
			Username: "admin",
		})

		if err != nil {
			log.Printf("Couldn't create json response because %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(resp)
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func validateUsernameRequest(r *http.Request) (string, bool) {
	token := r.Header.Get(TokenRequestHeader)
	return token, token != ""
}

type LoginTokenRequest struct {
	Username string
	Password string
}

type LoginTokenResponse struct {
	Token string
}

type UsernameResponse struct {
	Username string
}

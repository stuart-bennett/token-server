package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type userStore map[string]string

type tokenStore map[string]string

type app struct {
	Users  userStore
	Tokens tokenStore
}

func (a app) NewToken(user string) string {
	token := "test-token"
	a.Tokens[token] = user
	return token
}

func (a app) authenticate(u string, p string) bool {
	storedPwd, ok := a.Users[u]
	if !ok {
		return false
	}

	return p == storedPwd
}

func (a app) verifyToken(t string) bool {
	_, ok := a.Tokens[t]
	return ok
}

const TokenRequestHeader string = "X-Auth-Token"
const LoginPath string = "/login"
const UsernamePath string = "/username"

func main() {
	log.Fatal(http.ListenAndServe("localhost:8000", ConfigureMux()))
}

func ConfigureMux() *http.ServeMux {
	a := newApp()
	m := http.NewServeMux()
	m.Handle(LoginPath, http.HandlerFunc(a.login))
	m.Handle(UsernamePath, http.HandlerFunc(a.username))
	return m
}

func newApp() app {
	ts := tokenStore{}

	// seed with an initial user
	us := userStore{
		"admin": "admin1000",
	}
	return app{Users: us, Tokens: ts}
}

func (a app) login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ltr, ok := validateLoginTokenRequest(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.authenticate(ltr.Username, ltr.Password) {
		// Add token to store, associated with username
		token := a.NewToken(ltr.Username)
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(LoginTokenResponse{
			Token: token,
		})

		if err != nil {
			if err != nil {
				log.Printf("Couldn't create json response because %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func newLoginToken() string {
	return "temp-token"
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

func (a app) username(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	token, ok := validateUsernameRequest(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.verifyToken(token) {
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

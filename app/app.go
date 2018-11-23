package app

import (
	"net/http"

	"github.com/stuart-bennett/token-server/urls"
)

const TokenRequestHeader string = "X-Auth-Token"

type TokenStore interface {
	NewToken(user string) string
	VerifyToken(token string) (string, bool)
}

type userStore map[string]string

type App struct {
	Users  userStore
	Tokens TokenStore
}

func New(ts TokenStore) App {
	// seed with an initial user
	us := userStore{
		"admin":        "admin1000",
		"someone-else": "password123",
	}

	return App{Users: us, Tokens: ts}
}

func (a App) Authenticate(u string, p string) bool {
	storedPwd, ok := a.Users[u]
	if !ok {
		return false
	}

	return p == storedPwd
}

func ConfigureMux(tokenStore TokenStore) *http.ServeMux {
	endpoints := New(tokenStore)
	m := http.NewServeMux()
	m.Handle(string(urls.Login), http.HandlerFunc(endpoints.Login))
	m.Handle(string(urls.Username), http.HandlerFunc(endpoints.Username))
	return m
}

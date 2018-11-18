package app

import (
	"github.com/stuart-bennett/token-server/stores"
)

const TokenRequestHeader string = "X-Auth-Token"

type userStore map[string]string

type App struct {
	Users  userStore
	Tokens stores.TokenStore
}

func NewApp(ts stores.TokenStore) App {
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

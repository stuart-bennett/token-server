package app

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
)

const TokenRequestHeader string = "X-Auth-Token"

type userStore map[string]string
type tokenStore map[string]string
type App struct {
	Users  userStore
	Tokens tokenStore
}

func NewApp() App {
	ts := tokenStore{}

	// seed with an initial user
	us := userStore{
		"admin": "admin1000",
	}

	return App{Users: us, Tokens: ts}
}

func (a App) NewToken(user string) string {
	token := newToken()
	a.Tokens[token] = user
	return token
}

func newToken() string {
	// A sequence of 100 random numbers
	n := 100
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("Could not generate a secure random byte sequence to create a login token because %s", err))
	}

	checksum := sha512.Sum512(b)
	return fmt.Sprintf("%x", checksum)
}

func (a App) Authenticate(u string, p string) bool {
	storedPwd, ok := a.Users[u]
	if !ok {
		return false
	}

	return p == storedPwd
}

func (a App) VerifyToken(t string) bool {
	_, ok := a.Tokens[t]
	return ok
}

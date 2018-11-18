package stores

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
)

type TokenStore interface {
	NewToken(user string) string
	VerifyToken(token string) bool
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

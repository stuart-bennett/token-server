package tokenstore

import (
	"crypto/rand"
	"fmt"
)

func newToken() string {
	// A sequence of 100 random numbers
	n := 100
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("Could not generate a secure random byte sequence to create a login token because %s", err))
	}

	return fmt.Sprintf("%x", b)
}

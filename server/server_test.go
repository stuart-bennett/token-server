package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stuart-bennett/token-server/app"
	"github.com/stuart-bennett/token-server/testhelper"
	"github.com/stuart-bennett/token-server/tokenstore"
	"github.com/stuart-bennett/token-server/urls"
)

var ts *httptest.Server
var c http.Client

func TestMain(t *testing.T) {
	setup()
	defer teardown()
	t.Run("Login endpoint", func(t *testing.T) {
		t.Run("Should only accept POST requests", func(t *testing.T) {
			testhelper.OnlyAcceptsMethod(
				http.MethodPost,
				urls.Login.Abs(ts),
				t, &c)
		})

		t.Run(
			"Malformed request should respond with 400",
			func(t *testing.T) { testMalformedRequest(t, ts) })

		t.Run(
			"Invalid credentials should respond with 401",
			func(t *testing.T) { testInvalidCredentials(t, ts) })

		t.Run(
			"When successful should return JSON response containing login token",
			func(t *testing.T) { testLoginResponse(t, ts) })
	})

	t.Run("Username endpoint", func(t *testing.T) {
		t.Run("Should only accept GET requests", func(t *testing.T) {
			testhelper.OnlyAcceptsMethod(
				http.MethodGet,
				urls.Username.Abs(ts),
				t, &c)
		})

		t.Run(
			"Missing token header value should respond with 401",
			func(t *testing.T) { testMissingAuthHeader(t, ts) })

		t.Run(
			"Invalid token in header should respond with 401",
			func(t *testing.T) { testInvalidToken(t, ts) })

		t.Run(
			"Valid token in header should respond with 200 & username in body",
			func(t *testing.T) { testValidToken(t, ts) })
	})
}

func setup() {
	ts = httptest.NewServer(app.ConfigureMux(tokenstore.InMemory{}))
	c = http.Client{}
}

func teardown() {
	ts.Close()
}

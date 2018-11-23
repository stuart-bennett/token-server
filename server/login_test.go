package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stuart-bennett/token-server/app"
	"github.com/stuart-bennett/token-server/testhelper"
	"github.com/stuart-bennett/token-server/urls"
)

func testMalformedRequest(t *testing.T, ts *httptest.Server) {
	reqs := []string{
		"{ \"username\": 100, \"password\": true }",
		"{ \"username\": \"100\", \"password\": true }",
		"{ \"username\": 100, \"password\": \"test\" }",
		"{ \"username\": \"100\", }",
		"{ \"password\": true }",
		"{ }",
	}

	for _, req := range reqs {
		t.Run(req, func(t *testing.T) {
			res, err := http.Post(
				urls.Login.Abs(ts),
				"application/json",
				bytes.NewBufferString(req))

			if err != nil {
				t.Error(err)
			}

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf(
					"%s - Got: %d. Want: %d",
					req,
					res.StatusCode,
					http.StatusBadRequest)
			}
		})
	}
}

func testInvalidCredentials(t *testing.T, ts *httptest.Server) {
	res, err := http.Post(
		urls.Login.Abs(ts),
		"application/json",
		testhelper.NewLoginRequestJson("admin", "admin999"))

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 401 {
		t.Errorf("Got: %d. Want %d", res.StatusCode, 401)
	}
}

func testLoginResponse(t *testing.T, ts *httptest.Server) {
	res, err := http.Post(
		urls.Login.Abs(ts),
		"application/json",
		testhelper.NewLoginRequestJson("admin", "admin1000"))

	if err != nil {
		t.Error(err)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Got: %s Want: %s", ct, "application/json")
	}

	var r = app.LoginTokenResponse{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Error(err)
	}

	if r.Token == "" {
		t.Fail()
	}
}

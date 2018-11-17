package main

import (
	"bytes"
	"encoding/json"
	"github.com/stuart-bennett/token-server/app"
	"github.com/stuart-bennett/token-server/testhelper"
	"net/http"
	"net/http/httptest"
	"testing"
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
		res, err := http.Post(
			ts.URL+LoginPath,
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
	}
}

func testInvalidCredentials(t *testing.T, ts *httptest.Server) {
	res, err := http.Post(
		ts.URL+LoginPath,
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
		ts.URL+LoginPath,
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

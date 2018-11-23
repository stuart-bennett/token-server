package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stuart-bennett/token-server/app"
	"github.com/stuart-bennett/token-server/testhelper"
	"github.com/stuart-bennett/token-server/urls"
)

func testMissingAuthHeader(t *testing.T, ts *httptest.Server) {
	// no X-Auth-Token header!
	res, err := http.Get(urls.Username.Abs(ts))
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Got: %d. Want: %d", res.StatusCode, http.StatusBadRequest)
	}
}

func testInvalidToken(t *testing.T, ts *httptest.Server) {
	req, err := http.NewRequest(http.MethodGet, urls.Username.Abs(ts), nil)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set(app.TokenRequestHeader, "not-valid")
	res, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Got: %d. Want: %d", res.StatusCode, http.StatusUnauthorized)
	}
}

func acquireValidToken() (string, error) {
	res, err := http.Post(
		urls.Login.Abs(ts),
		"application/json",
		testhelper.NewLoginRequestJson("admin", "admin1000"))

	if err != nil {
		return "", err
	}

	var rsp app.LoginTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&rsp); err != nil {
		return "", err
	}

	if rsp.Token == "" {
		return "", fmt.Errorf("No token was produced")
	}

	return rsp.Token, nil
}

func testValidToken(t *testing.T, ts *httptest.Server) {
	req, err := http.NewRequest(http.MethodGet, urls.Username.Abs(ts), nil)
	if err != nil {
		t.Error(err)
	}

	token, err := acquireValidToken()
	if err != nil {
		t.Error(err)
	}

	req.Header.Set(app.TokenRequestHeader, token)
	res, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Got: %d. Want: %d", res.StatusCode, http.StatusOK)
	}

	var resp app.UsernameResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Error(err)
	}

	if resp.Username != "admin" {
		t.Fail()
	}
}

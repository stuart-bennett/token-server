package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ts *httptest.Server
var c http.Client

func TestMain(t *testing.T) {
	setup()
	defer teardown()
	t.Run("/login should only accept POST requests", func(t *testing.T) {
		onlySingleMethod(http.MethodPost, LoginPath, t)
	})

	t.Run("Username should only accept GET requests", func(t *testing.T) {
		onlySingleMethod(http.MethodGet, UsernamePath, t)
	})

	t.Run("When using valid credentials, POST /login should return JSON response containing login token", LoginResponseTest)
}

func LoginResponseTest(t *testing.T) {
	res, err := http.Post(ts.URL + LoginPath, "application/json", nil)
	if err != nil {
		t.Error(err)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Got: %s Want: %s", ct, "application/json")
	}

	var r = LoginTokenResponse{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Error(err)
	}

	if r.Token == "" {
		t.Fail()
	}
}

func setup() {
	ts = httptest.NewServer(ConfigureMux())
	c = http.Client{}
}

func teardown() {
	ts.Close()
}

func filter(items []string, f func(item string) bool) []string {
	var result []string
	for _, item := range items {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}

var standardHttpMethods [8]string = [8]string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

func onlySingleMethod(m string, path string, t *testing.T) {
	ms := filter(standardHttpMethods[0:], func(item string) bool {
		return item != m
	})
	target := ts.URL + path
	for _, m := range ms {
		req, err := http.NewRequest(m, target, nil)
		if err != nil {
			t.Error(err)
		}

		res, err := c.Do(req)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Method = %s, Response Status: %d", m, res.StatusCode)
		}
	}

	req, err := http.NewRequest(m, target, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Method = %s, Response Status: %d", http.MethodPost, res.StatusCode)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ts *httptest.Server
var c http.Client

func TestMain(t *testing.T) {
	setup()
	defer teardown()
	t.Run("Login endpoint", func(t *testing.T) {
		t.Run("Should only accept POST requests", func(t *testing.T) {
			onlySingleMethod(http.MethodPost, LoginPath, t)
		})

		t.Run(
			"Malformed request should respond with 400",
			testMalformedRequest)

		t.Run(
			"Invalid credentials should respond with 401",
			invalidCredentialsTest)

		t.Run(
			"When successful should return JSON response containing login token",
			loginResponseTest)
	})

	t.Run("/username should only accept GET requests", func(t *testing.T) {
		onlySingleMethod(http.MethodGet, UsernamePath, t)
	})
}

func testMalformedRequest(t *testing.T) {
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

func invalidCredentialsTest(t *testing.T) {
	res, err := http.Post(
		ts.URL+LoginPath,
		"application/json",
		newLoginRequestJson("admin", "admin999"))

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 401 {
		t.Errorf("Got: %d. Want %d", res.StatusCode, 401)
	}
}

func newLoginRequestJson(u string, p string) *bytes.Buffer {
	json, err := json.Marshal(LoginTokenRequest{
		Username: u,
		Password: p,
	})

	if err != nil {
		panic(fmt.Sprintf("Couldn't create login request, %s", err))
	}

	return bytes.NewBuffer(json)
}

func loginResponseTest(t *testing.T) {
	res, err := http.Post(
		ts.URL+LoginPath,
		"application/json",
		newLoginRequestJson("admin", "admin1000"))

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

func (ms StandardHttpMethods) filter(f func(item string) bool) []string {
	var result []string
	for _, item := range ms {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}

type StandardHttpMethods [8]string

var standardHttpMethods StandardHttpMethods = StandardHttpMethods{
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
	ms := standardHttpMethods.filter(func(item string) bool {
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

	if res.StatusCode == http.StatusMethodNotAllowed {
		t.Errorf("Method = %s, Response Status: %d", http.MethodPost, res.StatusCode)
	}
}

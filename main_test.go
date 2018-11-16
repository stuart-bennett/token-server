package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var ts *httptest.Server
var c http.Client

func TestMain(t *testing.T) {
	setup()
	defer teardown()
	t.Run("Login should only accept POST requests", OnlyAcceptsPostTest)
	t.Run("Username should only accept GET requests", OnlyAcceptsGetTest)
	teardown()
}

func setup() {
	ts = httptest.NewServer(ConfigureMux())
	c = http.Client{}
}

func teardown() {
	ts.Close()
}

func OnlyAcceptsPostTest(t *testing.T) {
	ms := [3]string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodPatch,
	}

	target := ts.URL + LoginPath
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

	req, err := http.NewRequest(http.MethodPost, target, nil)
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

func OnlyAcceptsGetTest(t *testing.T) {
	ms := [3]string{
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
	}

	target := ts.URL + UsernamePath
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

	req, err := http.NewRequest(http.MethodGet, target, nil)
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

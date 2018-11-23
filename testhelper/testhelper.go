package testhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stuart-bennett/token-server/app"
)

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

func (ms StandardHttpMethods) filter(f func(item string) bool) []string {
	var result []string
	for _, item := range ms {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}

func NewLoginRequestJson(u string, p string) *bytes.Buffer {
	json, err := json.Marshal(app.LoginTokenRequest{
		Username: u,
		Password: p,
	})

	if err != nil {
		panic(fmt.Sprintf("Couldn't create login request, %s", err))
	}

	return bytes.NewBuffer(json)
}

func OnlyAcceptsMethod(m string, target string, t *testing.T, c *http.Client) {
	ms := standardHttpMethods.filter(func(item string) bool {
		return item != m
	})

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

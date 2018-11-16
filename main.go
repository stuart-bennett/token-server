package main

import (
	"log"
	"net/http"
)

type tokenStore map[string]string

const LoginPath string = "/login"
const UsernamePath string = "/username"

func main() {
	log.Fatal(http.ListenAndServe("localhost:8000", ConfigureMux()))
}

func ConfigureMux() *http.ServeMux {
	ts := tokenStore{}
	m := http.NewServeMux()
	m.Handle(LoginPath, http.HandlerFunc(ts.login))
	m.Handle(UsernamePath, http.HandlerFunc(ts.username))
	return m
}

func (ts tokenStore) login(w http.ResponseWriter, req *http.Request) {
}

func (ts tokenStore) username(w http.ResponseWriter, req *http.Request) {
}

package main

import (
	"github.com/stuart-bennett/token-server/app"
	"log"
	"net/http"
)

const (
	LoginPath string = "/login"
	UsernamePath string = "/username"
)

func main() {
	log.Fatal(http.ListenAndServe("localhost:8000", ConfigureMux()))
}

func ConfigureMux() *http.ServeMux {
	endpoints := app.NewApp()
	m := http.NewServeMux()
	m.Handle(LoginPath, http.HandlerFunc(endpoints.Login))
	m.Handle(UsernamePath, http.HandlerFunc(endpoints.Username))
	return m
}

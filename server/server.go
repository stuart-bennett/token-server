package main

import (
	"flag"
	"github.com/stuart-bennett/token-server/app"
	"github.com/stuart-bennett/token-server/stores"
	"log"
	"net/http"
	"strconv"
)

const (
	LoginPath    string = "/login"
	UsernamePath string = "/username"
	ListenPort   int    = 8000
)

var redisAddr = flag.String(
	"r", "",
	"IP address and port, e.g. 172.17.0.2:6367. Activates the Redis-based token store and uses the instance at the address provided")

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(
		":"+strconv.Itoa(ListenPort),
		ConfigureMux(getStore())))
}

func ConfigureMux(tokenStore stores.TokenStore) *http.ServeMux {
	endpoints := app.NewApp(tokenStore)
	m := http.NewServeMux()
	m.Handle(LoginPath, http.HandlerFunc(endpoints.Login))
	m.Handle(UsernamePath, http.HandlerFunc(endpoints.Username))
	return m
}

func getStore() stores.TokenStore {
	if *redisAddr != "" {
		log.Printf("Using Redis Token Store @ %s", *redisAddr)
		return stores.NewRedisTokenStore(*redisAddr)
	}

	log.Printf("Using in-memory token store")
	return stores.InMemoryTokenStore{}
}

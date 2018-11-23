package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/stuart-bennett/token-server/app"
	"github.com/stuart-bennett/token-server/tokenstore"
)

const ListenPort int = 8000

var redisAddr = flag.String(
	"r", "",
	"IP address and port, e.g. 172.17.0.2:6367. Activates the Redis-based token store and uses the instance at the address provided")

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(
		":"+strconv.Itoa(ListenPort),
		app.ConfigureMux(getStore())))
}

func getStore() app.TokenStore {
	if *redisAddr != "" {
		log.Printf("Using Redis Token Store @ %s", *redisAddr)
		return tokenstore.NewRedis(*redisAddr)
	}

	log.Printf("Using in-memory token store")
	return tokenstore.InMemory{}
}

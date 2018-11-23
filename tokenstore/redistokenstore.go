package tokenstore

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

const TokenKeyPrefix string = "token"

type Redis struct {
	Pool *redis.Pool
}

func NewRedis(addr string) Redis {
	log.Printf("[RedisTokenStore] Connecting to %s", addr)
	pool := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}

	return Redis{Pool: &pool}
}

func (ts Redis) NewToken(user string) string {
	token := newToken()
	conn := ts.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", makeKeyName(token), user)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func (ts Redis) VerifyToken(token string) (string, bool) {
	log.Printf(token)
	conn := ts.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", makeKeyName(token)))
	if err != nil {
		log.Printf("Problem verifying token: %s", err)
		return "", false
	}

	if !exists {
		return "", false
	}

	username, err := redis.String(conn.Do("GET", makeKeyName(token)))
	if err != nil {
		log.Printf("Problem verifying token: %s", err)
		return "", false
	}

	return username, true
}

func makeKeyName(s string) string {
	return TokenKeyPrefix + ":" + s
}

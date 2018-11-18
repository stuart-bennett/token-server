package stores

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

const TokenKeyPrefix string = "token"

type RedisTokenStore struct {
	Pool *redis.Pool
}

func NewRedisTokenStore(addr string) RedisTokenStore {
	pool := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}

	return RedisTokenStore{Pool: &pool}
}

func (ts RedisTokenStore) NewToken(user string) string {
	token := newToken()
	conn := ts.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", TokenKeyPrefix+":"+token, user)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func (ts RedisTokenStore) VerifyToken(token string) bool {
	log.Printf(token)
	conn := ts.Pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", TokenKeyPrefix+":"+token))
	if err != nil {
		log.Printf("Problem verifying token: %s", err)
		return false
	}

	if !exists {
		return false
	}

	value, err := redis.String(conn.Do("GET", TokenKeyPrefix+":"+token))
	if err != nil {
		log.Printf("Problem verifying token: %s", err)
		return false
	}

	log.Printf("Found %s", value)
	return true
}

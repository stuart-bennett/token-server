package tokenstore

import "sync"

type InMemory map[string]string

var mu = sync.RWMutex{}

func (ts InMemory) NewToken(user string) string {
	token := newToken()
	mu.Lock()
	ts[token] = user
	mu.Unlock()
	return token
}

func (ts InMemory) VerifyToken(token string) (string, bool) {
	mu.RLock()
	username, ok := ts[token]
	mu.RUnlock()
	return username, ok
}

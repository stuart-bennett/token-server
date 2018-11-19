package stores

import "sync"

type InMemoryTokenStore map[string]string

var mu = sync.RWMutex{}

func (ts InMemoryTokenStore) NewToken(user string) string {
	token := newToken()
	mu.Lock()
	ts[token] = user
	mu.Unlock()
	return token
}

func (ts InMemoryTokenStore) VerifyToken(token string) (string, bool) {
	mu.RLock()
	username, ok := ts[token]
	mu.RUnlock()
	return username, ok
}

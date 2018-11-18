package stores

type InMemoryTokenStore map[string]string

func (ts InMemoryTokenStore) NewToken(user string) string {
	token := newToken()
	ts[token] = user
	return token
}

func (ts InMemoryTokenStore) VerifyToken(token string) bool {
	_, ok := ts[token]
	return ok
}

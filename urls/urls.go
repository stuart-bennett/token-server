package urls

import "net/http/httptest"

type Url string

const (
	Login    Url = "/login"
	Username Url = "/username"
)

func (u Url) String() string {
	return string(u)
}

func (u Url) Abs(s *httptest.Server) string {
	return s.URL + string(u)
}

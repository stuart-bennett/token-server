package app

import (
	"encoding/json"
	"log"
	"net/http"
)

type UsernameResponse struct {
	Username string `json:"username"`
}

// GET /username
func (a App) Username(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token, ok := validateUsernameRequest(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, ok := a.Tokens.VerifyToken(token)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(UsernameResponse{
		Username: username,
	})

	if err != nil {
		log.Printf("Couldn't create json response because %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func validateUsernameRequest(r *http.Request) (string, bool) {
	token := r.Header.Get(TokenRequestHeader)
	return token, token != ""
}

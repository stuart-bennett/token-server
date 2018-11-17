package app

import (
	"encoding/json"
	"log"
	"net/http"
)

type UsernameResponse struct {
	Username string
}

// GET /username
func (a App) Username(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	token, ok := validateUsernameRequest(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.VerifyToken(token) {
		resp, err := json.Marshal(UsernameResponse{
			Username: "admin",
		})

		if err != nil {
			log.Printf("Couldn't create json response because %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(resp)
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func validateUsernameRequest(r *http.Request) (string, bool) {
	token := r.Header.Get(TokenRequestHeader)
	return token, token != ""
}

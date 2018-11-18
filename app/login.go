package app

import (
	"encoding/json"
	"log"
	"net/http"
)

type LoginTokenRequest struct {
	Username string
	Password string
}

type LoginTokenResponse struct {
	Token string
}

// POST /login
func (a App) Login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ltr, ok := validateLoginTokenRequest(req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if a.Authenticate(ltr.Username, ltr.Password) {
		// Add token to store, associated with username
		token := a.Tokens.NewToken(ltr.Username)
		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(LoginTokenResponse{
			Token: token,
		})

		if err != nil {
			if err != nil {
				log.Printf("Couldn't create json response because %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func validateLoginTokenRequest(req *http.Request) (LoginTokenRequest, bool) {
	var ltr LoginTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&ltr); err != nil {
		log.Printf("Could not decode request because %s", err)
		return LoginTokenRequest{}, false
	}

	if ltr.Username == "" || ltr.Password == "" {
		return LoginTokenRequest{}, false
	}

	return ltr, true
}

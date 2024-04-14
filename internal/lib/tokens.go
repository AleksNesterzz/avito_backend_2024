package tokens

import (
	apierr "avito_backend/internal/lib/errors"
	"net/http"
)

var (
	AdminToken  = "admin"
	CasualToken = "user"
)

func GetToken(w http.ResponseWriter, r *http.Request) (string, error) {
	token := r.Header.Get("token")
	if token == "" {
		return "", apierr.ErrNoAuth
	} else {
		return token, nil
	}
}

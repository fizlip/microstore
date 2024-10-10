package middleware

import (
	"errors"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

var UnAuthorizedError = errors.New("Invalid username or token.")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var provider string = r.URL.Query().Get("provider")
		var token = r.Header.Get("Authorization")
		var err error

		if provider == "" || token == "" {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		err = services.ValidateAccessToken(token, provider)

		if err != nil {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}
        log.Println("User is authenticated")
		next.ServeHTTP(w, r)
	})
}

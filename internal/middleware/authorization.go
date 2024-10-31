package middleware

import (
	"errors"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

var UnAuthorizedError = errors.New("Forbidden resource")
var AuthenticationError = errors.New("Missing credentials")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token = r.Header.Get("Authorization")
        var provider = r.Header.Get("X-Provider")
		var instigator = r.Header.Get("X-Username")
		var err error
        
		if provider == "" || token == "" || instigator == ""{
			log.Error(AuthenticationError)
			api.AuthenticationErrorHandler(w, AuthenticationError)
			return
		}
		_, err = services.ValidateAccessToken(token, provider, instigator)

		if err != nil {
			log.Error(UnAuthorizedError)
			api.AuthorizationErrorHandler(w, UnAuthorizedError)
			return
		}
        log.Println("User is authenticated")
		next.ServeHTTP(w, r)
	})
}

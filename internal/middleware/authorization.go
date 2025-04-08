// package middleware contains functions that process responses 
// before they are handled 
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

// Authorization will return an error if the user is not correctly authenticated.
// A user is correctly authenticated if they have a valid username and a valid 
// token. Each token needs to be accompanied by a provider. Currently we
// allow users to authenticat using GitHub and Google.
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
        // This function will return an error if the OAuth token given 
        // by the provider is not valid 
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

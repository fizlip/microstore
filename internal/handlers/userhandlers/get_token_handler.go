package userhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

// GetToken TODO: Docs
func GetToken(w http.ResponseWriter, r *http.Request) {
	var err error

	customer_id := r.Header.Get("X-Username")
	server := r.URL.Query().Get("server")
	tokenService, err := services.NewTokenService()
	if err != nil {
		log.Error("Cannot connect to customer service", err)
		api.InternalErrorHandler(w)
		return
	}

	token, err := tokenService.GetToken(customer_id, server)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	// And run the product with it
	// result, err := (*user).GetInvoice()
	var code = 200
	if len(token.Token.EncryptedData) == 0 {
		code = 404
	}
	err = json.NewEncoder(w).Encode(api.TokenSuccess{Code: code, Token: token.Token})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

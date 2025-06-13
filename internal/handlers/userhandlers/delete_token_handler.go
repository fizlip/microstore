package userhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

// GetToken TODO: Docs
func DeleteToken(w http.ResponseWriter, r *http.Request) {
	var err error

	customer_id := r.Header.Get("X-Username")
	server := r.URL.Query().Get("server")
	tokenService, err := services.NewTokenService()
	if err != nil {
		log.Error("Cannot connect to customer service", err)
		api.InternalErrorHandler(w)
		return
	}

	err = tokenService.DeleteToken(customer_id, server)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	// And run the product with it
	// result, err := (*user).GetInvoice()
	err = json.NewEncoder(w).Encode(api.TokenSuccess{Code: 200})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

package userhandlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/models"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

// GetToken TODO: Docs
func CreateToken(w http.ResponseWriter, r *http.Request) {
	var err error

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
	var token models.Token
	err = json.Unmarshal(body, &token)

	tokenService, err := services.NewTokenService()
	if err != nil {
		log.Error("Cannot connect to customer service", err)
		api.InternalErrorHandler(w)
		return
	}

	err = tokenService.CreateToken(token)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	// And run the product with it
	// result, err := (*user).GetInvoice()
	err = json.NewEncoder(w).Encode(api.TokenSuccess{Code: 200, Token: ""})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

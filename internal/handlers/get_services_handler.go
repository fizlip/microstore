package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/tools"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetUserServices(w http.ResponseWriter, r *http.Request) {
	var params = api.UserServiceParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	var tokenDetails *tools.ServiceDetails
	tokenDetails = (*database).GetUserServices(params.Username)
	if tokenDetails == nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var response = api.UserServiceResponse{
		Balance: (*tokenDetails).Service,
		Code:    http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

}

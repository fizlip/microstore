// package vendorhandlers contains request handlers for all endpoints that
// interact with the vendor
package storehandlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/models"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// GetTools will return all tools owned by the given vendor
// example: https://<hostname>/get/tools
// response:
//
//	SUCCESS: 200, <slice of tools>
//	FAILURE: 500, An unexpected error occured
func CreateCommunication(w http.ResponseWriter, r *http.Request) {
	var err error

	var store *services.Store
	store, err = services.NewStoreController()
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	var comm models.Communication
	err = json.Unmarshal(body, &comm)

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = store.CreateCommunication(comm)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = json.NewEncoder(w).Encode(api.Success{Code: 200, Message: "Communication created"})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

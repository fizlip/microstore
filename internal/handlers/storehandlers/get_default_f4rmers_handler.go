// package vendorhandlers contains request handlers for all endpoints that
// interact with the vendor
package storehandlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// GetTools will return all tools owned by the given vendor
// example: https://<hostname>/get/tools
// response:
//
//	SUCCESS: 200, <slice of tools>
//	FAILURE: 500, An unexpected error occured
func GetDefaultF4rmers(w http.ResponseWriter, r *http.Request) {
	var err error
	var user *services.Customer

	user, err = services.NewCustomer("default")
	result, err := user.ReadF4rmers()
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

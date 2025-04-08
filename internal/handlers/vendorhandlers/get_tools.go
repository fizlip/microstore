// package vendorhandlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetTools will return all tools owned by the given vendor
// example: https://<hostname>/get/tools
// response: 
//      SUCCESS: 200, <slice of tools> 
//      FAILURE: 500, An unexpected error occured
func GetTools(w http.ResponseWriter, r *http.Request) {
    vendor_id := r.Header.Get("X-Username")
    vendor, err := services.NewVendorController(vendor_id)
	if err != nil {
		log.Fatalf("failed to retrieve usage data, %v", err)
        return
	}

    result, err := vendor.GetTools()
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }

	err = json.NewEncoder(w).Encode(api.ActionsSuccess{Code: 200, Message: result})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

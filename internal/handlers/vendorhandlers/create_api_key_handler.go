// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
    "net/http"
    "fmt"
    "encoding/json"

    "github.com/fizlip/microstore/internal/services"
    "github.com/fizlip/microstore/api"

	log "github.com/sirupsen/logrus"
)

// CreateAPIKey will generate a new API key regiestered to a vendor.
// To create a new API key a vendor needs to be logged in and create a request.
//
// example: POST https://<hostname>/vendor/key/create
// response:
//      SUCCESS: 200, <api_key>
//      FAILURE: 500, An unexepcted error occurred
func CreateAPIKey(w http.ResponseWriter, r *http.Request) {
    var err error
    var vendor *services.Vendor
    var username = r.Header.Get("X-Username")

    vendor, err = services.NewVendorController(username)

    apiKey, err := vendor.CreateAPIKey(username)
    if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
    
    response := api.Success{
        Code: 200, 
        Message: fmt.Sprintf("%s", apiKey),
    }

    err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

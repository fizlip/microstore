// package f4rmer_handlers provides functions used to interact with f4rmers.
package apihealth

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// IsLive is used to see whether the api is online or not. Usually used as a utility for 
// devs when they want to know if the API is online but can also be used by 
// customers on the website to see if our services are online.
//
// example: GET https://<hostname>/health/
// response: 
//      SUCCESS: 200, "f4rmhouse api is live", 
//      FAILURE: 403, An unexpected error occured 
//
func IsLive(w http.ResponseWriter, r *http.Request) {
    var err error

    w.Header().Set("Content-Type", "application/json")
    res := api.Success{200, "f4rmhouse api is live!"}
    err = json.NewEncoder(w).Encode(res)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

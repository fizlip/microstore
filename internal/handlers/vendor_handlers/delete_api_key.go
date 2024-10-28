package vendor_handlers

import (
    "net/http"
    "fmt"
    "encoding/json"

    "github.com/fizlip/microstore/internal/services"
    "github.com/fizlip/microstore/api"

	log "github.com/sirupsen/logrus"
)

func DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
    var err error
    var vendor *services.Vendor
    
    key := r.Header.Get("Authorization") 
    username := r.Header.Get("X-Username") 

    vendor, err = services.NewVendorController(username)

    err = vendor.DeleteAPIKey(username, key)
    if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
    
    response := api.Success{
        Code: 200, 
        Message: fmt.Sprintf("Deleted api key sucessfully"),
    }

    err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

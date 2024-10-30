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
    
    username    := r.Header.Get("X-Username") 
    partial     := r.Header.Get("X-Delete") 

    vendor, err = services.NewVendorController(username)

    keys, err := vendor.APIKeyStorage.Read(username)
    var sk string
    for _, key := range keys {
        if key.Partial == partial {
            sk = key.Key
        }
    }
    
    if sk == "" {
        log.Error("No such key")
		api.InternalErrorHandler(w)
		return
    }

    err = vendor.DeleteAPIKey(username, sk)
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

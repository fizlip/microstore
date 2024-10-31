package vendor_handlers

import (
    "net/http"
    "fmt"
    "encoding/json"

    "github.com/fizlip/microstore/internal/services"
    "github.com/fizlip/microstore/api"

	log "github.com/sirupsen/logrus"
)

func CreateAPIKey(w http.ResponseWriter, r *http.Request) {
    var err error
    var vendor *services.Vendor
    vendor, err = services.NewVendorController("fillatino@gmail.com")

    apiKey, err := vendor.CreateAPIKey("fillatino@gmail.com")
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

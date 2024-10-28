package vendor_handlers

import (
    "net/http"
    "encoding/json"

    "github.com/fizlip/microstore/internal/services"
    "github.com/fizlip/microstore/api"

	log "github.com/sirupsen/logrus"
)

func ReadAPIKey(w http.ResponseWriter, r *http.Request) {
    var err error
    var vendor *services.Vendor
    vendor, err = services.NewVendorController("fillatino@gmail.com")
    if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

    keys, err := vendor.ReadAllAPIKeys("fillatino@gmail.com")
    if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

    err = json.NewEncoder(w).Encode(keys)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}


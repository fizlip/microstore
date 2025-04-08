package storehandlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetUTIProduct(w http.ResponseWriter, r *http.Request) {
	var err error

    var store *services.Store
    store, err = services.NewStoreController()
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

    uti := r.URL.Query().Get("uti")

	result, err := (*store).GetUTIProduct(uti)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    var code int = 200
    if result.UID == "" {
       code = 201 
    }

    err = json.NewEncoder(w).Encode(api.ActionSuccess{Code: code, Message: result})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

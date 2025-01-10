package storehandlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	var err error

    var store *services.Store
    store, err = services.NewStoreController()
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

    uid := r.URL.Query().Get("uid")

	result, err := (*store).GetProduct(uid)
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

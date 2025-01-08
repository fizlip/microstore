package storehandlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetTrending will return a slice of trending products.
//
// example: GET https://<hostname>/store/trending
// response: 
//      SUCCESS: 200, <slice of products>
//      FAILURE: 500, An unexpected error occurred
func GetTrending(w http.ResponseWriter, r *http.Request) {
	var err error

    var store *services.Store
    store, err = services.NewStoreController()

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	result, err := (*store).GetTrendingBusiness()
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

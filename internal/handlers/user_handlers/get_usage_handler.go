package user_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

// GetUsage returns a list of tools that have been used by the user 
// and how much they have been used
func GetUsage(w http.ResponseWriter, r *http.Request) {
	var err error

    customer_id := r.URL.Query().Get("uid")
	user, err := services.NewCustomer(customer_id)
	if err != nil {
        log.Error("Cannot connect to customer service", err)
        api.InternalErrorHandler(w)
        return
	}

	// And run the product with it
	result, err := (*user).GetAnalytics()
    if err != nil {
        log.Error("Cannot get analytics from database", err)
        api.InternalErrorHandler(w)
        return
    }

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

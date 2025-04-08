package userhandlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

    "net/http"

	log "github.com/sirupsen/logrus"
)

// CreateEmptyRecord will create empty records for the users analytics 
func CreateEmptyRecord(w http.ResponseWriter, r *http.Request) {

	var err error

    w.Header().Set("Content-Type", "application/json")

    customer_id := r.URL.Query().Get("uid")
    user, err := services.NewCustomer(customer_id)
	if err != nil {
        log.Error("Cannot connect to customer service", err)
        api.InternalErrorHandler(w)
        return
	}

    err = (*user).CreateEmptyAnalyticsRecord() 
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    } 

    res := api.Success{200, "Empty analytics record created"}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

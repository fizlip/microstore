package userhandlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

    "net/http"

	log "github.com/sirupsen/logrus"
)

// CheckValidSession lets a client check whether they have valid credentials. Returns an error 
// iff the user is unauthorized otherwise returns 200. 
func CreateEmptyRecord(w http.ResponseWriter, r *http.Request) {

	var err error
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

	err = json.NewEncoder(w).Encode("New user record created")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

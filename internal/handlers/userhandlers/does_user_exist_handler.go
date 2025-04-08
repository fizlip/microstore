package userhandlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
    // "github.com/fizlip/microstore/internal/models"

    "net/http"

	log "github.com/sirupsen/logrus"
)

// DoesUserExist will return true if a user with a given uid is in the database
// this is used to check whether this is the first time a user logs in.
func DoesUserExist(w http.ResponseWriter, r *http.Request) {

	var err error
    var res bool = false

    customer_id := r.Header.Get("X-Username")
    if customer_id == "" {
        log.Error("No customer_id specified")
        api.InternalErrorHandler(w)
        return
    }

    user, err := services.NewCustomer(customer_id)
	if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
	}

    user_analytics, err := (*user).GetAnalytics() 
    if len(user_analytics) == 1 {res = true}

	err = json.NewEncoder(w).Encode(api.ExistsSuccess{Code: 200, Exists:res})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

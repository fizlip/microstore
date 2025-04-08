package userhandlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

    "net/http"

	log "github.com/sirupsen/logrus"
)

// CreateEmptyRecord will create empty records for the users analytics 
func HasUsedAction(w http.ResponseWriter, r *http.Request) {

	var err error
    var username = r.Header.Get("X-Username")
    uti := r.URL.Query().Get("uti")

    w.Header().Set("Content-Type", "application/json")


    user, err := services.NewCustomer(username)
	if err != nil {
        log.Error("Cannot connect to customer service", err)
        api.InternalErrorHandler(w)
        return
	}

    res, err := (*user).HasUsedAction(uti) 
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    } 

	err = json.NewEncoder(w).Encode(api.ExistsSuccess{Code: 200, Message: res})
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

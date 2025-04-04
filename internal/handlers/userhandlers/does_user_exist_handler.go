package userhandlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
    // "github.com/fizlip/microstore/internal/models"

    "net/http"

	log "github.com/sirupsen/logrus"
)

// CheckValidSession lets a client check whether they have valid credentials. Returns an error 
// iff the user is unauthorized otherwise returns 200. 
func DoesUserExist(w http.ResponseWriter, r *http.Request) {

	var err error
    var res bool = false
    customer_id := r.URL.Query().Get("uid")
    user, err := services.NewCustomer(customer_id)
	if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
	}

    user_analytics, err := (*user).GetAnalytics() 
    if len(user_analytics) == 1 {res = true}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

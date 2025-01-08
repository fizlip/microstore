package userhandlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "net/http"

	log "github.com/sirupsen/logrus"
)

// CheckValidSession lets a client check whether they have valid credentials. Returns an error 
// iff the user is unauthorized otherwise returns 200. 
func CheckValidSession(w http.ResponseWriter, r *http.Request) {

	var err error
	w.Header().Set("Content-Type", "application/json")

    res := api.Success{200, "Session is valid"}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

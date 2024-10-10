package handlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "net/http"

	log "github.com/sirupsen/logrus"
)

// CheckValidSession is a simple endpoint that lets a client see whether they
// have valid credentials it returns an error iff the user is unauthorized
// otherwise it returns a simple message. 
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

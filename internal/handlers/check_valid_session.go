package handlers

import (
    "encoding/json"
    "github.com/fizlip/microstore/api"
    "net/http"

	log "github.com/sirupsen/logrus"
)

// GetAllImages will return one page of all images available in the f4mhouse OCI repo
// A correct response will return an array of Image values.
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

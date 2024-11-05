// Package review_handlers contains the request handlers for all 
// endpoints that are interact with reviews.
package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// UpdateReview will update a row in the Reviews table in the db
// TODO: is this needed?
func UpdateReview(w http.ResponseWriter, r *http.Request) {
	var err error

	err = json.NewEncoder(w).Encode("")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

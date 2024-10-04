package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// DeleteReview will remove a row in the Reviews table in the db
// ex: <hostname>/user/review/create
// data: ...
func DeleteReview(w http.ResponseWriter, r *http.Request) {
	var err error

	log.Println("Delete a review")

	err = json.NewEncoder(w).Encode("")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

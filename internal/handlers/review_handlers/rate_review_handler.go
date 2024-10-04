package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// RateReview will ...
// ex: <hostname>/user/review/create
// data: ...
func RateReview(w http.ResponseWriter, r *http.Request) {
	var err error

	log.Println("Rate a review")

	err = json.NewEncoder(w).Encode("")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// UpdateReview will update a row in the Reviews table in the db
// ex: <hostname>/user/review/create
// data: ...
func UpdateReview(w http.ResponseWriter, r *http.Request) {
	var err error

    reviewService, err := services.NewReviewService()
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }

    // TODO: need to add read of query param somewhere
    err = reviewService.Delete("123")
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }
    
	err = json.NewEncoder(w).Encode("")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

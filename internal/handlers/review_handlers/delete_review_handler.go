// Package review_handlers contains the request handlers for all 
// endpoints that are interact with reviews.
package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// DeleteReview will remove a row in the Reviews table in the db
// ex: <hostname>/user/review/delete?uid=<REVIEW_ID>
//
// response: 500 if error 200 if OK 
func DeleteReview(w http.ResponseWriter, r *http.Request) {
	var err error
    // Create service
    var reviews *services.ReviewService
	reviews, err = services.NewReviewService()
	if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
        return 
	}

    review_id := r.URL.Query().Get("uid")

    // Execute function
    err = (*reviews).Delete(review_id)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
        return 
	}

    response := api.Success{
        Code:    200,
	    Message: "Review was successfully deleted",
    }

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

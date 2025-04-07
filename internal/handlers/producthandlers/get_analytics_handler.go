// Package reviewhandlers contains the request handlers for all 
// endpoints that are interact with reviews.
package producthandlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// DeleteReview will remove a row in the Reviews table in the db
// example: https://<hostname>/user/review/delete?uid=<REVIEW_ID>
//
// response: 
//      SUCCESS: 200, Review successfully deleted
//      FAILURE: 500, An unexpected error occurred
func GetActionAnalytics(w http.ResponseWriter, r *http.Request) {
	var err error

    uti := r.URL.Query().Get("uti")

    // Create service
	analytics, err := services.NewActionAnalyticsService()
	if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
        return 
	}

    // Execute function
    result, err := analytics.GetRecentCalls(uti)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
        return 
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

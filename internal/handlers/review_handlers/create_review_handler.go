// Package review_handlers contains the request handlers for all 
// endpoints that are interact with reviews.
package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/internal/models"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// CreateReview will insert another row into the Reviews table in the db
//
// example: https://<hostname>/user/review/create
// body: { 
//  PID: string
//  Title: string
//  Value: string
//  Content: string
//  Creator: string
//  Date: string
// }
//
// response: 
//      SUCCESS: 200, Review was successfully created 
//      FAILURE: 500, An unexpected error occured
func CreateReview(w http.ResponseWriter, r *http.Request) {
    var response api.Success
    var err error
    var reviews *services.ReviewService
    reviews, err = services.NewReviewService()
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return 
    }
    
    // Read body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return 
    }

    var review models.Review
    err = json.Unmarshal(body, &review)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    } else {
        err = (*reviews).Create(review)

        response = api.Success{
            Code:    200,
            Message: "Review was successfully created",
        }
    }

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

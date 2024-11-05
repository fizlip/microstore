// Package review_handlers contains the request handlers for all 
// endpoints that are interact with reviews.
package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetReviews will return a slice of all reviews connected to a given product_id 
//
// Since each review has an optional interaction (like/no like) we also check to see whether
// the instigator of the request has rated the review. If yes the Rated field 
// in the review is set to true, otherwise it's false. 
//
// ex: <hostname>/user/review/get?product_id=<PRODUCT_ID>
// data: []Review
// response: 500 if error, 200 if OK
func GetReviews(w http.ResponseWriter, r *http.Request) {
	var err error
	var reviewService *services.ReviewService
    var username = r.Header.Get("X-Username")

    queryParams := r.URL.Query()
    // get product id 
    pid := queryParams.Get("product_id")

	reviewService, err = services.NewReviewService()

	if err != nil {
		log.Fatalf("failed to get ECR images, %v", err)
	}

	reviews, err := (*reviewService).GetReviews(pid)
    // Get the UIDs of all reviews so that we can see if user has interacted with them
    var interactionUids []string
    for _, r := range reviews {
        interactionUids = append(interactionUids, r.UID)
    }
    
    // Check whether user has interacted with any review connected to the product
    // interactions will be an array of all interactions user has on this product
    interactions, err := reviewService.GetInteractions(interactionUids, username)   
    if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
    
    // if there are no interactions just return all reviews with rated = false 
    if len(interactions) == 0 {
	    err = json.NewEncoder(w).Encode(reviews)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
    } else {
        // Here we see which reviews user has rated
        // All reviews that user has rated will set Rated to true 
        for i, r := range reviews {
            for _, interaction := range interactions {
                if interaction.UID == r.UID {
                    reviews[i].Rated = true
                }
            }
        }
        err = json.NewEncoder(w).Encode(reviews)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
    }
}

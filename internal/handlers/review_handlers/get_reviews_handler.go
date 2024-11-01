package review_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetReviews ...
// ex: <hostname>/user/review/create
// data: ...
func GetReviews(w http.ResponseWriter, r *http.Request) {
	var err error
	var reviews *services.ReviewService
    var username = r.Header.Get("X-Username")

    queryParams := r.URL.Query()
    // get product id 
    pid := queryParams.Get("product_id")

	reviews, err = services.NewReviewService()

	if err != nil {
		log.Fatalf("failed to get ECR images, %v", err)
	}

	rs, err := (*reviews).GetReviews(pid)
    var interactionUids []string
    for _, r := range rs {
        interactionUids = append(interactionUids, r.UID)
    }
    
    is, err := reviews.GetInteractions(interactionUids, username)   
    if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
    
    if len(is) == 0 {
	    err = json.NewEncoder(w).Encode(rs)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
    } else {
        for i, r := range rs {
            for _, interaction := range is {
                if interaction.UID == r.UID {
                    rs[i].Rated = true
                }
            }
        }
        err = json.NewEncoder(w).Encode(rs)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
    }
}

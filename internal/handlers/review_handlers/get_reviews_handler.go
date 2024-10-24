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
	reviews, err = services.NewReviewService()

	if err != nil {
		log.Fatalf("failed to get ECR images, %v", err)
	}

	rs, err := (*reviews).GetReviews("123")

	err = json.NewEncoder(w).Encode(rs)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

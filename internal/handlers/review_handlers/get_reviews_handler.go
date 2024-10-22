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
	var reviewsTable *services.ITable
	reviewsTable, err = services.NewReviewTable("reviews")

	if err != nil {
		log.Fatalf("failed to get ECR images, %v", err)
	}

	rs, err := (*reviewsTable).Read("123")

	err = json.NewEncoder(w).Encode(rs)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

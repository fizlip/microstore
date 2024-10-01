package vendor_handlers 

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

// GetCurrent...
// ex: <hostname>/user/review/create
// data: ...
func GetCurrentUsage(w http.ResponseWriter, r *http.Request) {
  var err error
  var reviewsTable *api.ITable
  reviewsTable, err = api.NewReviewTable("reviews")

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

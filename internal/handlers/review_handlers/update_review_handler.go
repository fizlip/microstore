package review_handlers 

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

// UpdateReview will update a row in the Reviews table in the db
// ex: <hostname>/user/review/create
// data: ...
func UpdateReview(w http.ResponseWriter, r *http.Request) {
  var err error

  log.Println("Update a review")

  err = json.NewEncoder(w).Encode("")
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

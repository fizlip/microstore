package review_handlers 

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

// CreateReview will insert another row into the Reviews table in the db
// ex: <hostname>/user/review/create
// data: ...
func CreateReview(w http.ResponseWriter, r *http.Request) {
  var err error

  log.Println("Create a new review")

  err = json.NewEncoder(w).Encode("")
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

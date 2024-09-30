package review_handlers 

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/api"
  "io/ioutil"

  log "github.com/sirupsen/logrus"
)

// CreateReview will insert another row into the Reviews table in the db
// ex: <hostname>/user/review/create
// data: ...
func CreateReview(w http.ResponseWriter, r *http.Request) {
  var response api.Success
  var err error
  var reviewTable *api.ITable
  reviewTable, err = api.NewReviewTable("table")
  if err != nil {
    panic(err)
  }

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  var review api.Review
  err = json.Unmarshal(body, &review)
  if err != nil {
    response = api.Success {
      Code: 200,
      Message: "Review was successfully created",
    }
  } else {
    _, err = (*reviewTable).Insert(review)

    response = api.Success {
      Code: 200,
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

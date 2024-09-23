package handlers

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/internal/tools"
  "github.com/fizlip/microstore/api"
  "github.com/gorilla/schema"

  log "github.com/sirupsen/logrus"
)
// GetProduct will return one image with the UID given in the URL param
func GetProduct (w http.ResponseWriter, r *http.Request) {
  var params = api.GetProductParams{}
  var decoder *schema.Decoder = schema.NewDecoder()
  var err error

  err = decoder.Decode(&params, r.URL.Query())

  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }

  var database *tools.DatabaseInterface
  database, err = tools.NewDatabase()

  if err != nil {
    api.InternalErrorHandler(w)
    return
  }

  var response = (*database).GetProduct(params.UID)

  log.Println("Get product")

  w.Header().Set("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(response)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

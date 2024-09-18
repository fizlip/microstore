package handlers

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/internal/tools"
  "github.com/fizlip/microstore/api"
  "github.com/gorilla/schema"

  log "github.com/sirupsen/logrus"
)
// GetAllImages will return one page of all images available in the f4mhouse OCI repo
// A correct response will return an array of Image values.
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

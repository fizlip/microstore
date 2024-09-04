package handlers

import (
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/fizlip/microstore/internal/tools"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

// UpdateImage field, new_value
func UpdateImage(w http.ResponseWriter, r *http.Request) {
  var err error

  var database *tools.DatabaseInterface
  database, err = tools.NewDatabase()

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  var upd api.UpdateParams
  err = json.Unmarshal(body, &upd)
  if err != nil {
    panic(err)
  }

  // Call update method in db interface
  err = (*database).UpdateMetadata(upd)

  response := api.Success{
    Code: 200,
    Message: "Record was successfully updated",
  }

  err = json.NewEncoder(w).Encode(response)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

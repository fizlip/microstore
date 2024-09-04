package handlers

import (
  "fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/fizlip/microstore/internal/tools"
  "github.com/fizlip/microstore/api"
  "github.com/google/uuid"

  log "github.com/sirupsen/logrus"
)

func CreateImage(w http.ResponseWriter, r *http.Request) {
  var err error
  var database *tools.DatabaseInterface

  database, err = tools.NewDatabase()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }
  var meta tools.Metadata 
  err = json.Unmarshal(body, &meta)
  if err != nil {
    panic(err)
  }  

  if err != nil {
    api.InternalErrorHandler(w)
    return
  }

  meta.UUID = uuid.New().String()

  err = (*database).CreateMetadata(meta)

  w.Header().Set("Content-Type", "application/json")

  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
  response := api.Success{
    Code: 200,
    Message: fmt.Sprintf("New record with ID: %s has successfully been created", meta.UUID),
  }
  err = json.NewEncoder(w).Encode(response)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

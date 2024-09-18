package handlers

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/internal/tools"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

func GetTrending(w http.ResponseWriter, r *http.Request) {
  var err error
  var database *tools.DatabaseInterface
  database, err = tools.NewDatabase()

  if err != nil {
    api.InternalErrorHandler(w)
    return
  }
  
  result := (*database).GetTrending()

  w.Header().Set("Content-Type", "application/json")
  log.Println("Get trending")

  err = json.NewEncoder(w).Encode(result)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

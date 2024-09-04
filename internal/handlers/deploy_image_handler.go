package handlers

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

func DeployImage(w http.ResponseWriter, r *http.Request) {
  var err error
  err = json.NewEncoder(w).Encode("{}")
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

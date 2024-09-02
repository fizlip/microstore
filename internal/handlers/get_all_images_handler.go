package handlers

import (
  "encoding/json"
  "net/http"

  "github.com/fizlip/microstore/api"
  log "github.com/sirupsen/logrus"
)

func GetAllImages(w http.ResponseWriter, r *http.Request) {

  var err error

  var response = api.UserServiceResponse{
    Balance: 100,
    Code: http.StatusOK,
  }

  w.Header().Set("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(response)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

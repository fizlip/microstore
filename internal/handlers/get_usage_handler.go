package handlers

import (
  "encoding/json"
  "net/http"

  "github.com/fizlip/microstore/api"
  log "github.com/sirupsen/logrus"
)

func GetUsage(w http.ResponseWriter, r *http.Request) {
  var err error
  user, err := api.NewCustomer("test")
  if err != nil {
    panic(err)
  } 
    
  // And run the product with it
  result, err := (*user).GetAnalytics()
  err = json.NewEncoder(w).Encode(result)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}


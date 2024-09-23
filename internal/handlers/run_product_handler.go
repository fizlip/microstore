package handlers

import (
  //"fmt"
  "encoding/json"
  "net/http"

  "github.com/fizlip/microstore/api"
  log "github.com/sirupsen/logrus"
)

func RunProduct(w http.ResponseWriter, r *http.Request) {
  var err error
  user, err := api.NewCustomer("some_email")
  if err != nil {
    panic(err)
  } 
    
  // Create a meter event
  (*user).RunProduct()
  err = json.NewEncoder(w).Encode("")
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}


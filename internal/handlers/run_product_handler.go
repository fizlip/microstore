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
  evt := &api.F4SubscriptionEvent{
    CustomerID: "cus_QsQCZFMdRmYIP8", 
    Value: "10", 
    ProductID: "f4rmhouse_usage"}

  // And run the product with it
  res, err := (*user).RunProduct(evt)
  err = json.NewEncoder(w).Encode(res)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}


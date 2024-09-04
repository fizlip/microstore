package handlers

import (
  //"fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"

  "github.com/fizlip/microstore/api"
  log "github.com/sirupsen/logrus"
)

func CreateCluster(w http.ResponseWriter, r *http.Request) {
  var err error
  var ecsParams api.ECSParams

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  err = json.Unmarshal(body, &ecsParams)
  if err != nil {
    panic(err)
  }
  log.Println(ecsParams)

  response := api.Success {
    Code: 200,
    Message: "Cluster deployment has succseded",
  }

  err = json.NewEncoder(w).Encode(response)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

package handlers

import (
  //"fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"

  "github.com/fizlip/microstore/api"
  log "github.com/sirupsen/logrus"
)

func CreateTaskdef(w http.ResponseWriter, r *http.Request) {
  var err error
  var taskdef api.Taskdef

  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err)
  }

  err = json.Unmarshal(body, &taskdef)
  if err != nil {
    panic(err)
  }
  log.Println(taskdef)

  err = json.NewEncoder(w).Encode("{}")
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

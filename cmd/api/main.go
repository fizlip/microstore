package main

import (
  "fmt"
  "net/http"

  "github.com/go-chi/chi"
  "github.com/fizlip/microstore/internal/handlers"
  log "github.com/sirupsen/logrus"
)

func main() {
  log.SetReportCaller(true)

  var r *chi.Mux = chi.NewRouter()

  handlers.Handle(r)

  fmt.Println("Starting GO API service on localhost:8000...")
  fmt.Println(`
                _                                              
               (_)                          _                    
         ____  _  ____  ____ ___      ___ _| |_ ___   ____ _____ 
         |    \| |/ ___)/ ___) _ \    /___|_   _) _ \ / ___) ___ |
         | | | | ( (___| |  | |_|    |___ | | || |_| | |   | ____|
         |_|_|_|_|\____)_|   \___/   (___/   \__)___/|_|   |_____)
  `)

  err := http.ListenAndServe("localhost:8000", r)

  if err != nil {
    log.Error(err)
  }
}

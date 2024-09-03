package api

import (
  "encoding/json"
  "net/http"
)

type UserServiceParams struct {
  Username string
}

type UserServiceResponse struct {
  Code int
  Balance int64
}

type Error struct {
  Code int
  Message string
}

type Image struct {
  Name 		string `json:name`
  Description 	string `json:description`
  Rating 	string `json:rating`
  ImageID 	string `json:image_id`
  Tags 		string `json:tags`
  Creator 	string `json:creator`
  Requests 	string `json:requests`
  Price 	string `json:price`
}

func writeError(w http.ResponseWriter, message string, code int) {

  resp := Error {
    Code: code,
    Message: message,
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)

  json.NewEncoder(w).Encode(resp)

}

var (
  RequestErrorHandler = func(w http.ResponseWriter, err error) {
    writeError(w, err.Error(), http.StatusBadRequest)
  }
  InternalErrorHandler = func(w http.ResponseWriter) {
    writeError(w, "An unexpected Error occured", http.StatusInternalServerError)
  }
)

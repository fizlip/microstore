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

type UpdateParams struct {
  ID    string `json:"id"`
  Field string `json:"field"`
  Value string `json:"value"`
}

type Success struct {
  Code int
  Message string
}

type Error struct {
  Code int
  Message string
}

// AI generated
type Metadata struct {
  Digest      string    `json:"digest"`     // ID of the image associated with the item
  Created     int 	`json:"created"`      // Creation time of the item
  Title       string    `json:"title"`        // Title of the item
  Description string    `json:"description"`  // Description of the item
  Reviews     int       `json:"reviews"`      // Number of reviews
  Tag         string    `json:"tag"`          // Tag associated with the item
  Creator     string    `json:"creator"`      // Creator of the item
  Requests    int       `json:"requests"`     // Number of requests for the item
  Pulls       int       `json:"pulls"`        // Number of pulls (e.g., for a Docker image)
  Price       float64   `json:"price"`        // Price of the item
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

package tools

import (
  log "github.com/sirupsen/logrus"
  "github.com/aws/aws-sdk-go-v2/service/ecr/types"
  "github.com/fizlip/microstore/api"

)

type LoginDetails struct {
  AuthToken string
  Username string
}

type ServiceDetails struct {
  Service int64
  Username string
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
  UUID        string    `json:"uuid"`
}

type DatabaseInterface interface {
  GetUserLoginDetails(username string) *LoginDetails
  GetUserServices(username string) *ServiceDetails
  GetMetadata(imageDigest []types.ImageIdentifier) []Metadata
  CreateMetadata(metadata Metadata) error
  UpdateMetadata(upd api.UpdateParams) error
  SetupDatabase() error
}

func NewDatabase() (*DatabaseInterface, error) {

  var database DatabaseInterface = &mockDB{}

  var err error = database.SetupDatabase()

  if err != nil {
    log.Error(err)
    return nil, err
  }

  return &database, nil

}

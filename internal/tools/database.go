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

type DatabaseInterface interface {
  GetUserLoginDetails(username string) *LoginDetails
  GetUserServices(username string) *ServiceDetails
  GetMetadata(imageDigest []types.ImageIdentifier) []api.Product
  CreateMetadata(metadata api.Product) error
  UpdateMetadata(upd api.UpdateParams) error
  GetTrending() [9]api.Product
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

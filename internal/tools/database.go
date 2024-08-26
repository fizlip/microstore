package tools

import (
  log "github.com/sirupsen/logrus"
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

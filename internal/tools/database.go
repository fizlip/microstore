package tools

import (
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/api/models"
	log "github.com/sirupsen/logrus"
)

type LoginDetails struct {
	AuthToken string
	Username  string
}

type ServiceDetails struct {
	Service  int64
	Username string
}

type DatabaseInterface interface {
	GetUserLoginDetails(username string) *LoginDetails
	GetUserServices(username string) *ServiceDetails
	GetProduct(uid string) models.Product
	SearchProduct(query string) []models.Product
	GetMetadata(imageDigest []types.ImageIdentifier) []models.Product
	CreateMetadata(metadata models.Product) error
	UpdateMetadata(upd api.UpdateParams) error
	GetTrending() [9]models.Product
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

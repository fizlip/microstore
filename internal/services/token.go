package services

import (
	// "errors"

	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fizlip/microstore/internal/models"

	log "github.com/sirupsen/logrus"
)

type TokenService struct {
	TokenStorage models.DatabaseReadWriter[models.Token]
}

func (s *TokenService) GetToken(username, provider string) (models.Token, error) {
	return s.TokenStorage.ReadSpecific(username, provider)
}

func (s *TokenService) CreateToken(token models.Token) error {
	return s.TokenStorage.Create(token, tokenToAttributeMapConverter)
}

func (s *TokenService) DeleteToken(username, server string) error {
	return s.TokenStorage.DeleteDetailed(username, server)
}

func NewTokenService() (*TokenService, error) {
	tokenT, err := models.NewDynamoTable[models.Token]("f4-secrets", "username", "server")
	if err != nil {
		return &TokenService{}, err
	}
	tokenService := TokenService{TokenStorage: tokenT}
	return &tokenService, nil
}

func tokenToAttributeMapConverter(token models.Token) map[string]types.AttributeValue {
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02 15:04:05")

	tokenSecret, err := attributevalue.Marshal(token.Token)
	if err != nil {
		log.Error(err)
	}

	return map[string]types.AttributeValue{
		"username": &types.AttributeValueMemberS{Value: token.Username},
		"server":   &types.AttributeValueMemberS{Value: token.Server},
		"token":    tokenSecret,
		"created":  &types.AttributeValueMemberS{Value: formattedDate},
		"provider": &types.AttributeValueMemberS{Value: token.Provider},
	}
}

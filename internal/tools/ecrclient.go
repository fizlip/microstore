package tools

import (
  log "github.com/sirupsen/logrus"
  "github.com/aws/aws-sdk-go-v2/service/ecr"
  "github.com/aws/aws-sdk-go-v2/aws"

)

type ECRClientInterface interface {
  GetImages(input *ecr.ListImagesInput) (*ecr.ListImagesOutput, error)
}

func NewEcrClient(cfg aws.Config) (*ECRClientInterface, error) {
  var ecrClient ECRClientInterface = &mockECRClient{}
  log.Println("Mock ecr created")
  return &ecrClient, nil
}

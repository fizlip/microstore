package tools

import (
  log "github.com/sirupsen/logrus"
  "github.com/aws/aws-sdk-go-v2/service/ecr"
  "github.com/fizlip/microstore/api"

)

type ComputeClientInterface interface {
  GetImages(input *ecr.ListImagesInput) (*ecr.ListImagesOutput, error)
  CreateContainer() error
  CreateConfig(cfg api.Taskdef) error
  DeployCluster(param api.ClusterParams) error
  Connect() (interface{}, error)
}

func NewComputeClient() (*ComputeClientInterface, error) {
  var computeClient ComputeClientInterface = &mockComputeClient{}
  log.Println("Mock ecr created")
  return &computeClient, nil
}

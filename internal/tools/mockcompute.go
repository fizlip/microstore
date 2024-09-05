package tools

import (
  "context"
  "github.com/aws/aws-sdk-go-v2/service/ecr"
  "github.com/aws/aws-sdk-go-v2/service/ecr/types"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/fizlip/microstore/api"
  "github.com/aws/aws-sdk-go-v2/config"

  "time"
  "math/rand"

  log "github.com/sirupsen/logrus"

)

func generateRandomString(length int) string {
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

// generatePlaceholderImageIds creates an array of 5 placeholder ImageIdentifier structs with random values
func generatePlaceholderImageIds(amt int) []types.ImageIdentifier {
  imageIds := make([]types.ImageIdentifier, amt)
  for i := 0; i < amt; i++ {
    imageIds[i] = types.ImageIdentifier{
      ImageDigest: aws.String(generateRandomString(64)), // Simulating a random 64-character digest
      ImageTag:    aws.String(generateRandomString(10)), // Simulating a random 10-character tag
    }
  }
  return imageIds
}

var nextToken string = "nextToken"

type mockComputeClient struct{
  cfg aws.Config
}

var mockImages = ecr.ListImagesOutput{
  ImageIds: generatePlaceholderImageIds(10),
  NextToken: &nextToken,
}

func (e *mockComputeClient) Connect() (interface{}, error) {
  cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("default"))
  if err != nil {
    log.Fatalf("unable to load SDK config, %v", err)
    return nil, err 
  }

  e.cfg = cfg

  return cfg, nil
}

func (e *mockComputeClient) GetImages(input *ecr.ListImagesInput) (*ecr.ListImagesOutput, error) {
  log.Println("Fetching images...")
  log.Println("Current region: ", e.cfg.Region)

  return &mockImages, nil
}

func (e *mockComputeClient) CreateContainer() error {
  log.Println("Creating container, beep boop beep")
  log.Println("Current region: ", e.cfg.Region)

  return nil
}

func (e *mockComputeClient) CreateConfig(cfg api.Taskdef) error {
  log.Println("Creating config...")
  log.Println("Current region: ", e.cfg.Region)

  return nil
}

func (e *mockComputeClient) DeployCluster(params api.ClusterParams) error {
  log.Println("Cluster deployment...")
  log.Println("Current region: ", e.cfg.Region)

  return nil
}

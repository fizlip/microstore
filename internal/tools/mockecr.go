package tools

import (
  "github.com/aws/aws-sdk-go-v2/service/ecr"
  "github.com/aws/aws-sdk-go-v2/service/ecr/types"
  "github.com/aws/aws-sdk-go-v2/aws"
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

type mockECRClient struct{}

var mockImages = ecr.ListImagesOutput{
  ImageIds: generatePlaceholderImageIds(10),
  NextToken: &nextToken,
}

func (e *mockECRClient) GetImages(input *ecr.ListImagesInput) (*ecr.ListImagesOutput, error) {
  log.Println("Fetching images...")

  return &mockImages, nil

}

package handlers

import (
  "encoding/json"
  "context"
  "net/http"
  "github.com/aws/aws-sdk-go-v2/service/ecr"
  "github.com/aws/aws-sdk-go-v2/config"
  "github.com/fizlip/microstore/internal/tools"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)
// GetAllImages will return one page of all images available in the f4mhouse OCI repo
// A correct response will return an array of Image values.
func GetAllImages(w http.ResponseWriter, r *http.Request) {

  var err error
  var name string = "f4rmhouse"
  var ecrClient *tools.ECRClientInterface
  var database *tools.DatabaseInterface
  database, err = tools.NewDatabase()

  if err != nil {
    api.InternalErrorHandler(w)
    return
  }

  // First fetch the OCI data
  cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("default"))
  if err != nil {
    log.Fatalf("unable to load SDK config, %v", err)
  }

  ecrClient, err = tools.NewEcrClient(cfg)
  listImagesInput := ecr.ListImagesInput{RepositoryName: &name}
  listImagesOutput, err := (*ecrClient).GetImages(&listImagesInput)

  if err != nil {
    log.Fatalf("failed to get ECR images, %v", err)
  }

  if len(listImagesOutput.ImageIds) == 0 {
    log.Warn("No OCIs found in ECR respository with name ", name)
  }

   

  // Then fetch the meta data
  result := (*database).GetMetadata(listImagesOutput.ImageIds)

  w.Header().Set("Content-Type", "application/json")

  err = json.NewEncoder(w).Encode(result)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

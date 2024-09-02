package handlers

import (
  "context"
  "fmt"
  "encoding/json"
  "net/http"
  "github.com/aws/aws-sdk-go-v2/aws"
  "github.com/aws/aws-sdk-go-v2/service/s3"
  "github.com/aws/aws-sdk-go-v2/config"

  "github.com/fizlip/microstore/api"
  log "github.com/sirupsen/logrus"
)

func GetAllImages(w http.ResponseWriter, r *http.Request) {

  var err error

  cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("default"))
  if err != nil {
    log.Fatalf("unable to load SDK config, %v", err)
  }

  svc := s3.NewFromConfig(cfg)
  result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

  if err != nil {
    log.Fatalf("unable to list buckets %v", err)
  }

  fmt.Println("Buckets:")
  for _, bucket := range result.Buckets {
    fmt.Printf("%s\n", aws.ToString(bucket.Name))
  }

  var response = api.UserServiceResponse{
    Balance: 100,
    Code: http.StatusOK,
  }

  w.Header().Set("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(response)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

package storehandlers

import (
    "context"
	"encoding/json"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetShowcase(w http.ResponseWriter, r *http.Request) {
	var err error

    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        api.InternalErrorHandler(w)
		return
    }

    bucketName := "f4-public"

    // Create an S3 client
    client := s3.NewFromConfig(cfg)

    bucket := services.NewBucketManager(client, bucketName)
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

    uti := r.URL.Query().Get("uti")

	result, err := (*bucket).ListFolder(context.TODO(), "showcases/" + uti)
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

    log.Println("Return: ", result)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

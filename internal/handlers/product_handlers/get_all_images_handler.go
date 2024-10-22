package product_handlers

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/tools"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetAllImages will return one page of all images available in the f4mhouse OCI repo
// A correct response will return an array of Image values.
func GetAllImages(w http.ResponseWriter, r *http.Request) {

	var err error
	var name string = "f4rmhouse"
	var computeClient *tools.ComputeClientInterface
	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	computeClient, err = tools.NewComputeClient()
	_, err = (*computeClient).Connect()
	listImagesInput := ecr.ListImagesInput{RepositoryName: &name}
	listImagesOutput, err := (*computeClient).GetImages(&listImagesInput)

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

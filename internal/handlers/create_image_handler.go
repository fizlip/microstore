package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/api/models"
	"github.com/fizlip/microstore/internal/tools"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// CreateImage will create the metadata for a given image.
func CreateImage(w http.ResponseWriter, r *http.Request) {
	var err error
	var database *tools.DatabaseInterface

	database, err = tools.NewDatabase()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var meta models.Product
	err = json.Unmarshal(body, &meta)
	if err != nil {
		panic(err)
	}

	// UUID is used to create a unique identifier just in case
	// two recrods are created with the same title at the same time.
	meta.UID = uuid.New().String()
	meta.RunEndpoint = fmt.Sprintf("https://run.botify.com/oci/%s", meta.UID)

	// Call create method in db interface
	err = (*database).CreateMetadata(meta)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	response := api.Success{
		Code:    200,
		Message: fmt.Sprintf("New record with ID: '%s' has successfully been created", meta.UID),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

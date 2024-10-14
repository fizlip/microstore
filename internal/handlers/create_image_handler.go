package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/api/models"
	"github.com/google/uuid"
	// "io/ioutil"
	"net/http"
    "errors"

	log "github.com/sirupsen/logrus"
)

var ToLargeError = errors.New("Form data to large")
var ImageNotFound = errors.New("Form does not contain image")

// CreateImage will create the metadata for a given image.
func CreateImage(w http.ResponseWriter, r *http.Request) {
	var err error
	// var database *tools.DatabaseInterface

    err = r.ParseMultipartForm(32 << 20) // 32MB max memory
    if err != nil {
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }

    file, fileHeader, err := r.FormFile("image")
    if err != nil {
        api.AuthenticationErrorHandler(w, ImageNotFound)
        return
    }
    log.Println("Filename: ", file)
    defer file.Close()

	// database, err = tools.NewDatabase()

	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	panic(err)
	//}

	var meta models.Product
	//err = json.Unmarshal(body, &meta)
	//if err != nil {
	//	panic(err)
	//}

    fname, err := services.CreateTempFile(file, fileHeader.Filename)
    if err != nil {
        api.AuthenticationErrorHandler(w, ImageNotFound)
        return
    }
    defer services.CleanupTempFile(fname)
	// UUID is used to create a unique identifier just in case
	// two recrods are created with the same title at the same time.
	meta.UID = uuid.New().String()
	// meta.RunEndpoint = fmt.Sprintf("https://run.botify.com/oci/%s", meta.UID)

	// Call create method in db interface
	// err = (*database).CreateMetadata(meta)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	response := api.Success{
		Code:    200,
		Message: fmt.Sprintf("New container with ID: '%s' has successfully been created", meta.UID),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

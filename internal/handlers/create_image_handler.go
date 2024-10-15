package handlers

import (
	"encoding/json"
	"fmt"

    "github.com/docker/docker/client"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/api/models"
	"github.com/google/uuid"
	// "io/ioutil"
	"net/http"
    "errors"
    "context"

	log "github.com/sirupsen/logrus"
)

var ToLargeError = errors.New("Form data to large")
var ImageNotFound = errors.New("Form does not contain image")
var DebugError = errors.New("Cannot list image")

// CreateImage will create the metadata for a given image.
func CreateImage(w http.ResponseWriter, r *http.Request) {
	var err error
	// var database *tools.DatabaseInterface

    // create a docker client object to communicate w/ docker daemon
    dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        panic(err)
    }

    _, err = dockerClient.Ping(context.Background())
    if err != nil {
        log.Error("Cannot connect to Docker daemon")
        api.AuthenticationErrorHandler(w, errors.New("Internal server error"))
        return
    }

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

    if err := services.DockerLoad(dockerClient, fname); err != nil {
        log.Println(err)
        api.AuthenticationErrorHandler(w, ImageNotFound)
        return
    }
    
    if err := services.PushToECR(dockerClient, "hello-world"); err != nil {
        log.Println(err)
        api.AuthenticationErrorHandler(w, ImageNotFound)
        return
    }

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

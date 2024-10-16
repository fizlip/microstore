package handlers

import (
	"encoding/json"
	"fmt"
    "strings"

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

// CreateImage will upload an image to the docker registry. It also creates the 
// database record used to identify the image 
// This endpoints needs a form header which should contain an image (which is a File)
// a title, descriptions, creator, price.
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

    imageFactory := services.DockerInterface{dockerClient}

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

    if err := imageFactory.DockerLoad(fname); err != nil {
        log.Println(err)
        api.AuthenticationErrorHandler(w, ImageNotFound)
        return
    }
    
    if err := imageFactory.PushToECR(strings.Split(fileHeader.Filename, ".")[0]); err != nil {
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

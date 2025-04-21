// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
    "context"
    "net/http"
    "fmt"
    "encoding/json"

    "github.com/docker/docker/client"
    "github.com/fizlip/microstore/internal/services"
    "github.com/fizlip/microstore/api"

	log "github.com/sirupsen/logrus"
)

// DeleteImage will ...
func DeleteImage(w http.ResponseWriter, r *http.Request) {
    var err error
    var vendor *services.Vendor
    var nameService *services.ToolNameService

    var username = r.Header.Get("X-Username")
    oci := r.URL.Query().Get("oci")

    // Init tools table this will initialize a connection with the dynamodb table
    vendor, err = services.NewVendorController(username)
    if err != nil {
        log.Error("Cannot connect to tools database ", err)
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }

    ownsTool := false
    isToolDeployed := false
    vendorTools, err := vendor.GetTools()
    if err != nil {
        log.Error("Cannot read tools")
        api.InternalErrorHandler(w)
        return
    }

    for _, tool := range vendorTools {
        if tool.UTI == oci && tool.Developer == username {
            ownsTool = true
            isToolDeployed = tool.Deployed
        }
    } 

    if !ownsTool {
        log.Error("Cannot delete image because user either does not own image or it does not exist")
        response := api.Success{
            Code: 201,
            Message: "User does not own tool or the tool does not exist",
        }
        err = json.NewEncoder(w).Encode(response)
	    if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
		    return
        }
        return
    }

    if isToolDeployed {
        log.Error("Cannot delete an image that is deployed")
        response := api.Success{
            Code: 202,
            Message: "Cannot delete an image that is deployed",
        }
        err = json.NewEncoder(w).Encode(response)
	    if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
		    return
        }
        return
    }

    // Create a docker client object to communicate w/ docker daemon
    dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        log.Error("Cannot create docker client")
        api.InternalErrorHandler(w)
        return
    }
    
    // Make sure docker daemon is running. Docker is used to push images to AWS ECR
    _, err = dockerClient.Ping(context.Background())
    if err != nil {
        log.Error("Cannot connect to Docker daemon")
        api.InternalErrorHandler(w)
        return
    }

    imageFactory := services.DockerInterface{dockerClient}

    log.Println("Deleting: ", oci)
    err = imageFactory.DeleteImageFromRepository("f4rmhouse", oci)

    if err != nil {
        log.Error("Cannot delete image: ", err)
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }

    // Init tools table this will initialize a connection with the dynamodb table
    nameService, err = services.NewToolNameService()
    if err != nil {
        log.Error("Cannot connect to TNS: ", err)
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }   

    err = nameService.DeleteRecord(oci)
    if err != nil {
        log.Error("Cannot delete name record from f4-ns: ", err)
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }

    err = vendor.DeleteTool(oci)
    if err != nil {
        log.Error("Cannot delete name record from f4-tools-prod: ", err)
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }

    err = json.NewEncoder(w).Encode(fmt.Sprintf("Image has been deleted: ", oci))
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

// package vendorhandlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
	"encoding/json"
    "net/http"

	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// DeployImage will create a cluster, taskdef and a new service for the cluster and taskdef.
// It should result in the image getting the correct hardware and being executable from an
// IP that has been given to the service by ECS.
//
// To deploy an image you need to have the unique tool identifier (UTI) which is the name of
// the tool as it has been registered in the tns (tool name service). 
//
// example: POST https://<hostname>/vendor/oci/deploy --form {
// {
//      tool_uti
// }
//
// response: 
//      SUCCESS: 200, Deployment finished without errors ... 
//      FAILURE: 500, An unexpected error occurred
func DeployImage(w http.ResponseWriter, r *http.Request) {
	var err error
    
    var deployer *services.DeployService 
    deployer, err = services.NewDeployService()

    var nameService *services.ToolNameService
    nameService, err = services.NewToolNameService()

    uti         := r.FormValue("tool_uti")
    taskDef     := uti + "-taskdef"
    serviceName := uti + "-service"

    err = deployer.CreateCluster(uti)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }
    err = deployer.CreateTaskDef(taskDef, "tool-test-1", "236798077409.dkr.ecr.eu-central-1.amazonaws.com/f4rmhouse:" + uti)
    err = deployer.CreateService(uti, serviceName, taskDef)

    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }

    if err := (*nameService).InsertIP(uti, "PENDING"); err != nil {
        log.Println("Could not insert tool into TNS", err)
        // TODO: add error for this case in /api
        api.CustomInternalErrorHandler(w, ConnectionError)
        return
    }
    
	response := api.Success{
		Code:    200,
		Message: "Deployment finished without errors. Use f4cli endpoints to get your tool URL. Note that it might take a couple of minutes to deploy hardware for your task.",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

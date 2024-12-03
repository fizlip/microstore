package vendor_handlers

import (
	"encoding/json"
    "net/http"

	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

func DeployImage(w http.ResponseWriter, r *http.Request) {
	var err error
    
    var deployer *services.DeployService 
    deployer, err = services.NewDeployService()

    uti         := r.FormValue("tool_uti")
    taskDef     := uti + "-taskdef"
    serviceName := uti + "-service"

    err = deployer.CreateCluster(uti)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }
    err = deployer.CreateTaskDef(taskDef, "tool-test-1", "236798077409.dkr.ecr.eu-central-1.amazonaws.com/f4rmhouse:grpc-server")
    err = deployer.CreateService(uti, serviceName, taskDef)

    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
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

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

    //deployer.CreateCluster("testNewCluster")
    //deployer.CreateCluster("testNewCluster")
    err = deployer.CreateTaskDef("tool-test-1", "236798077409.dkr.ecr.eu-central-1.amazonaws.com/f4rmhouse:grpc-server")

    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }
    
	log.Println("Deploy image?? (I dont know how to do this....)")

	response := api.Success{
		Code:    200,
		Message: "Deployment finished without errors. UTI: example.tool.com",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

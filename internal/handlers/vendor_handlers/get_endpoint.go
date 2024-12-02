package vendor_handlers

import (
	"encoding/json"
    "net/http"

	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

func GetEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
    
    var deployer *services.DeployService 
    deployer, err = services.NewDeployService()

    _, err = deployer.GetEndpoint("testNewCluster")

	response := api.Success{
		Code:    200,
		Message: "example.tool.com",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

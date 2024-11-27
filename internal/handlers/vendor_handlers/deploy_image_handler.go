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

    deployer.CreateCluster("testNewCluster")
    
	log.Println("Deploy image?? (I dont know how to do this....)")

	response := api.Success{
		Code:    200,
		Message: "Deployment in progress, you will recieve an email when it's finished! ",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

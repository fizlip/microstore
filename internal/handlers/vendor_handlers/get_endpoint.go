package vendor_handlers

import (
    "strings"
	"encoding/json"
    "net/http"

	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

func GetEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
    var message string
    
    var deployer *services.DeployService 
    deployer, err = services.NewDeployService()

    ips, err := deployer.GetEndpoint("testNewCluster")
    
    if len(ips) == 0 {
        message = "Cluster has no services associated with it."
    } else {
        message = strings.Join(ips, ",")
    }

	response := api.Success{
		Code:    200,
		Message: message,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

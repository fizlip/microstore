// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
package vendor_handlers

import (
    "strings"
	"encoding/json"
    "net/http"

	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// GetEndpoint will get the IP of the a tool given its uti (unique tool identifier)
//
// example: https://<hostname>/oci/get/endpoint
// response: 
//      SUCCESS: 200, <slice of ips>
//      FAILURE: 500, An unexpected error occurred
func GetEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
    var message string
    
    var deployer *services.DeployService 
    var ns *services.ToolNameService

    deployer, err = services.NewDeployService()
    ns, err = services.NewToolNameService()

    uti := r.URL.Query().Get("uti")

    ips, err := deployer.GetEndpoint(uti)
    if len(ips) == 0 {
        message = "Cluster has no services associated with it."
    } else {
        message = strings.Join(ips, ",")
        err = ns.InsertIP(uti, message)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
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

// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
    "net/http"
    "fmt"
    "encoding/json"

    "github.com/fizlip/microstore/internal/services"
    "github.com/fizlip/microstore/api"

	log "github.com/sirupsen/logrus"
)

// DeleteDeploymentHandler will ...
func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
    var err error
    var username = r.Header.Get("X-Username")

    var deployer *services.DeployService 
    deployer, err = services.NewDeployService()

    var vendor *services.Vendor
    vendor, err = services.NewVendorController(username)
    if err != nil {
        log.Error("Cannot connect to tools database ", err)
        api.AuthenticationErrorHandler(w, ToLargeError)
        return
    }

    uti := r.URL.Query().Get("uti")

    log.Println("Deleting: ", uti)

    err = deployer.DeleteService(uti)
    err = deployer.DeleteCluster(uti)
    err = deployer.DeleteTaskDef(uti)

    err = vendor.SetToolDeployed(uti, "false")
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }

    err = json.NewEncoder(w).Encode(fmt.Sprintf("Deployment has been deleted: ", uti))
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

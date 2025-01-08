// package vendorhandlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/internal/tools"
	"io/ioutil"
	"net/http"

	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

// CreateCluster
// TODO: this probably is not needed
//
//
func CreateCluster(w http.ResponseWriter, r *http.Request) {
	var err error
	var clusterParams api.ClusterParams
	var computeClient *tools.ComputeClientInterface

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &clusterParams)
	if err != nil {
		panic(err)
	}

	computeClient, err = tools.NewComputeClient()
	_, err = (*computeClient).Connect()

	err = (*computeClient).DeployCluster(clusterParams)

	response := api.Success{
		Code:    200,
		Message: "Cluster deployment has succseded",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

// package vendorhandlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
//
//	{
//	     tool_uti
//	}
//
// response:
//
//	SUCCESS: 200, Deployment finished without errors ...
//	FAILURE: 500, An unexpected error occurred
func DeployImage(w http.ResponseWriter, r *http.Request) {
	var err error
	var username = r.Header.Get("X-Username")

	var deployer *services.DeployService
	deployer, err = services.NewDeployService()

	var nameService *services.ToolNameService
	nameService, err = services.NewToolNameService()

	var vendor *services.Vendor
	vendor, err = services.NewVendorController(username)
	if err != nil {
		log.Error("Cannot connect to tools database ", err)
		api.AuthenticationErrorHandler(w, ToLargeError)
		return
	}

	var store *services.Store
	store, err = services.NewStoreController()
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	uti := r.FormValue("tool_uti")
	secretsStr := r.FormValue("secrets")
	taskDef := uti + "-taskdef"
	serviceName := uti + "-service"
	// [{"key":"AWS_ACCESS_KEY_ID","value":"AKIATOISD5HQ7CFF5PRK"},{"key":"AWS_REGION","value":"eu-central-1"},{"key":"AWS_SECRET_ACCESS_KEY","value":"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"}]
	// [{"key": ..., "value": ...}, ...]
	// [{"Key": ..., "Value": ...}, ...]

	log.Println("Secrets: ", secretsStr)

	secrets, err := parseKeyValuePairs(secretsStr)
	log.Println("Parsed Secrets: ", secrets)
	err = deployer.CreateCluster(uti)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
	err = deployer.CreateTaskDef(taskDef, "tool-test-1", "236798077409.dkr.ecr.eu-central-1.amazonaws.com/f4rmhouse:"+uti, secrets)
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

	err = store.VerifyProduct(uti)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = vendor.SetToolDeployed(uti, "true")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

type rawKV struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

func parseKeyValuePairs(jsonStr string) ([]types.KeyValuePair, error) {
	var rawPairs []rawKV
	if err := json.Unmarshal([]byte(jsonStr), &rawPairs); err != nil {
		return nil, err
	}

	var result []types.KeyValuePair
	for _, kv := range rawPairs {
		result = append(result, types.KeyValuePair{
			Name:  kv.Name,
			Value: kv.Value,
		})
	}

	return result, nil
}

package vendorhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// UpdateImageData will upload an image to the docker registry. It also creates the
// database record used to identify the image
// This endpoints needs a form header which should contain an image (which is a File)
// a title, descriptions, creator, price.
//
// example: POST https://<hostname>/vendor/oci/create --form
//
//	{
//	    title: string
//	    uti: string (unique tool identifier)
//	    description: string
//	    price: number
//	    thumbnail: string
//	    communityURL: string
//	    releaseType: string
//	    version: string
//	    showcase: string
//	    tags: string
//	    license: string
//	    endpoint_names: []string
//	    endpoint_descriptions: []string
//	    endpoint_parameters: [][]string
//	}
//
// response:
//
//	SUCCESS: 200, New container with ID: '%s' has successfully been created
//	FAILURE: 500, An unexpected error occurred
func UpdateImageData(w http.ResponseWriter, r *http.Request) {

	var err error
	var vendor *services.Vendor
	var store *services.Store

	var username = r.Header.Get("X-Username")

	/** Init **/

	// Init tools table this will initialize a connection with the dynamodb table
	vendor, err = services.NewVendorController(username)
	if err != nil {
		log.Error("Cannot connect to tools database ", err)
		api.AuthenticationErrorHandler(w, ToLargeError)
		return
	}

	// Init tools table this will initialize a connection with the dynamodb table
	store, err = services.NewStoreController()
	if err != nil {
		log.Error("Cannot connect to Store ", err)
		api.AuthenticationErrorHandler(w, ToLargeError)
		return
	}

	// Extract data from form values and create product struct based on them
	product, err := getProductFromForm(r)
	if err != nil {
		log.Error("Invalid form: ", err)
		api.AuthenticationErrorHandler(w, ToLargeError)
		return
	}

	oldTool, err := (*store).GetUTIProduct(product.UTI)
	if err != nil {
		log.Println("Name doesn't exist", err)
		api.AuthenticationErrorHandler(w, NameExistsError)
		return
	}

	product.UID = oldTool.UID
	product.Created = oldTool.Created
	product.Deployed = oldTool.Deployed
	log.Println("product", product)

	if err := (*vendor).CreateTool(product); err != nil {
		log.Println("Could not update product ", err)
		api.CustomInternalErrorHandler(w, ConnectionError)
		return
	}

	/** Result **/
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	// Success message
	sMsg := fmt.Sprintf("New container with ID: '%s' has successfully been created", product.UID)
	response := api.Success{
		Code:    200,
		Message: sMsg,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
	/** End result **/

}

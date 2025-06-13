package vendorhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/models"
	"github.com/fizlip/microstore/internal/services"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

var ToLargeError = errors.New("Error in form. Form data might be to large or your missing a file to upload.")
var ImageNotFound = errors.New("Form does not contain image")
var DebugError = errors.New("Cannot list image")
var ConnectionError = errors.New("Internal connection error. This is an issue on our side. We're working on it.")
var NameExistsError = errors.New("The name you gave your tool is already in use. Please pick another name")

// CreateImage will upload an image to the docker registry. It also creates the
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
func CreateImage(w http.ResponseWriter, r *http.Request) {

	var err error
	var vendor *services.Vendor
	var nameService *services.ToolNameService

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
	nameService, err = services.NewToolNameService()
	if err != nil {
		log.Error("Cannot connect to TNS ", err)
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

	_, err = (*nameService).GetTool(product.UTI)
	if err == nil {
		log.Println("Name already exists ", err)
		api.AuthenticationErrorHandler(w, NameExistsError)
		return
	}

	// Create a docker client object to communicate w/ docker daemon
	if product.Public {
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Error("Cannot create docker client")
			api.InternalErrorHandler(w)
			return
		}

		// Make sure docker daemon is running. Docker is used to push images to AWS ECR
		_, err = dockerClient.Ping(context.Background())
		if err != nil {
			log.Error("Cannot connect to Docker daemon")
			api.InternalErrorHandler(w)
			return
		}

		/** End init **/

		/** Execution **/
		// imageFactory is used to interface with the docker client.
		// DockerInterface exposes some utility functions used to load, tag and push
		// images
		imageFactory := services.DockerInterface{dockerClient}

		// We don't allow form data to be to large (this is not a restriction on image
		// size)
		err = r.ParseMultipartForm(32 << 20) // 32MB max memory
		if err != nil {
			api.PayloadTooLargeHandler(w, ToLargeError)
			return
		}

		// Get image file from form
		file, fileHeader, err := r.FormFile("image")
		if err != nil {
			api.CustomInternalErrorHandler(w, ImageNotFound)
			return
		}
		defer file.Close()

		// Create /tmp file that will be pushed
		fname, err := services.CreateTempFile(file, fileHeader.Filename)
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}
		defer services.CleanupTempFile(fname)

		// Loads image into docker
		if err := imageFactory.DockerLoad(fname); err != nil {
			log.Println("Docker image not found ", err)
			api.CustomInternalErrorHandler(w, ImageNotFound)
			return
		}

		// Tags and pushes image according to AWS ECR specifications
		if err := imageFactory.PushToECR(strings.Split(fileHeader.Filename, ".")[0]); err != nil {
			log.Println(err)
			api.CustomInternalErrorHandler(w, ConnectionError)
			return
		}
	}

	// Create a tools item in the DB
	if err := (*vendor).CreateTool(product); err != nil {
		log.Println("Could not create new product ", err)
		// TODO: add error for this case in /api
		api.CustomInternalErrorHandler(w, ConnectionError)
		return
	}

	if err := (*nameService).CreateRecord(
		models.ToolName{UTI: product.UTI, IPs: []string{}, Endpoints: product.Endpoints, Descriptions: product.EndpointDescriptions, Parameters: []string{}, DeploymentType: product.DeploymentType}); err != nil {
		log.Println("Could not insert tool into TNS", err)
		// TODO: add error for this case in /api
		api.CustomInternalErrorHandler(w, ConnectionError)
		return
	}

	/** End exeuction **/

	/** Result **/
	w.Header().Set("Content-Type", "application/json")
	// if err != nil {
	// 	log.Error(err)
	// 	api.InternalErrorHandler(w)
	// 	return
	// }

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

// getProductFromForm is a utility function used to extract fields from
// the form in the request and generates a models.Tool from them
func getProductFromForm(r *http.Request) (models.Tool, error) {
	// Extract form values
	title := r.FormValue("title")
	uti := r.FormValue("uti")
	description := r.FormValue("description")
	rating := 0.0
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	thumbnail := r.FormValue("thumbnail")
	overview := ""
	communityURL := r.FormValue("communityURL")
	reviews := 0
	developer := r.Header.Get("X-Username")
	pricingType := "pay-as-you-go"
	releaseType := r.FormValue("releaseType")
	version := r.FormValue("version")
	showcase := strings.Split(r.FormValue("showcase"), ",") // comma-separated values
	tags := strings.Split(r.FormValue("tags"), ",")         // comma-separated values
	runEndpoint := ""
	pulls := 0
	license := r.FormValue("license")
	verified := false
	endpoint_names := strings.Split(r.FormValue("endpoints_names"), ",")
	endpoint_desc := strings.Split(r.FormValue("endpoints_descriptions"), "<f4-delimiter>")
	endpoint_ps := r.FormValue("endpoints_parameters")
	deploymentType := strings.ToLower(r.FormValue("deployment_type"))
	public, parseErr := strconv.ParseBool(r.FormValue("public"))
	if parseErr != nil {
		log.Error("'public' value must be a bool", parseErr)
		return models.Tool{}, err
	}
	if err != nil {
		return models.Tool{}, err
	}

	log.Println("endpoints: ", endpoint_ps)

	var endpoints []models.Endpoint
	json.Unmarshal([]byte(endpoint_ps), &endpoints)

	// Generate error if required fields are empty
	if len(title) == 0 ||
		len(thumbnail) == 0 ||
		len(uti) == 0 ||
		len(developer) == 0 ||
		len(pricingType) == 0 ||
		len(releaseType) == 0 ||
		len(version) == 0 ||
		len(showcase) == 0 ||
		len(deploymentType) == 0 ||
		len(license) == 0 {

		log.Error("Missing value make sure that all product fields are present")
		return models.Tool{}, errors.New("Missing value in form")
	}

	// If no deployment type is selected or if the deployment type is not MCP or REST
	if !slices.Contains([]string{"mcp_sse", "mcp_streamable", "rest"}, deploymentType) {
		log.Error("deployment_type needs to be either 'MCP_SSE', 'MCP_STREAMABLE' or 'REST'. Current request has selected the invalid: ", deploymentType)
		return models.Tool{}, errors.New("Cannot deploy action without a valid deployment type.")
	}

	return models.Tool{
		UID:                  uuid.New().String(),
		Title:                title,
		UTI:                  uti,
		Description:          description,
		Rating:               rating,
		Price:                price,
		Thumbnail:            thumbnail,
		Overview:             overview,
		CommunityURL:         communityURL,
		Reviews:              reviews,
		Developer:            developer,
		PricingType:          pricingType,
		ReleaseType:          releaseType,
		Version:              version,
		Showcase:             showcase,
		Tags:                 tags,
		RunEndpoint:          runEndpoint,
		Pulls:                pulls,
		License:              license,
		Verified:             verified,
		EndpointNames:        endpoint_names,
		EndpointDescriptions: endpoint_desc,
		Endpoints:            endpoints,
		DeploymentType:       deploymentType,
		Public:               public,
	}, nil
}

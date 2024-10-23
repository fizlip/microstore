package vendor_handlers

import (
	"encoding/json"
	"fmt"
    "strings"
    "strconv"
    "net/http"
    "errors"
    "context"

    "github.com/docker/docker/client"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

var ToLargeError = errors.New("Error in form. Form data might be to large or your missing a file to upload.")
var ImageNotFound = errors.New("Form does not contain image")
var DebugError = errors.New("Cannot list image")
var ConnectionError = errors.New("Internal connection error. This is an issue on our side. We're working on it.")

// CreateImage will upload an image to the docker registry. It also creates the 
// database record used to identify the image 
// This endpoints needs a form header which should contain an image (which is a File)
// a title, descriptions, creator, price.
func CreateImage(w http.ResponseWriter, r *http.Request) {

	var err error
	var table *models.DynamoTable[models.Tool]

    /** Init **/

    // Init tools table this will initialize a connection with the dynamodb table
    table, err = models.NewDynamoTable[models.Tool]("f4-tools", "developer", "uid")
    if err != nil {
        log.Error("Cannot connect to tools database ", err)
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
    // Create a docker client object to communicate w/ docker daemon
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

    // Create a tools item in the DB
    if err := table.Create(product, convertToolToAWSItem); err != nil {
        log.Println("Could not create new product ", err)
        // TODO: add error for this case in /api
        api.CustomInternalErrorHandler(w, ConnectionError)
        return
    }
    
    /** End exeuction **/
    
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

func convertToolToAWSItem(tool models.Tool) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "uid": &types.AttributeValueMemberS{Value: tool.UID},
        "developer": &types.AttributeValueMemberS{Value: tool.Developer},
        "title": &types.AttributeValueMemberS{Value: tool.Title},
        "description": &types.AttributeValueMemberS{Value: tool.Description},
        "rating": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", tool.Rating)},
        "Pricing": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", tool.Price)},
        "thumbnail": &types.AttributeValueMemberS{Value: tool.Thumbnail},
        "overview": &types.AttributeValueMemberS{Value: tool.Overview},
        "communityURL": &types.AttributeValueMemberS{Value: tool.CommunityURL},
        "reviews": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", tool.Reviews)},
        "pricingType": &types.AttributeValueMemberS{Value: tool.PricingType},
        "releaseType": &types.AttributeValueMemberS{Value: tool.ReleaseType},
        "version": &types.AttributeValueMemberS{Value: tool.Version},
        "showcase": &types.AttributeValueMemberSS{Value: tool.Showcase},
        "tags": &types.AttributeValueMemberSS{Value: tool.Tags},
        "runEndpoint": &types.AttributeValueMemberS{Value: tool.RunEndpoint},
        "pulls": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", tool.Pulls)},
        "license": &types.AttributeValueMemberS{Value: tool.License},
        "verified": &types.AttributeValueMemberBOOL{Value: tool.Verified},
    }
}

// getProductFromForm is a utility function used to extract fields from
// the form in the request and generates a models.Tool from them
func getProductFromForm(r *http.Request) (models.Tool , error) {
    // Extract form values
    title           := r.FormValue("title")
    description     := r.FormValue("description")
    rating          := 0.0
    price, err      := strconv.ParseFloat(r.FormValue("price"), 64)
    thumbnail       := r.FormValue("thumbnail")
    overview        := ""
    communityURL    := r.FormValue("communityURL")
    reviews         := 0
    developer       := r.Header.Get("X-Username")
    pricingType     := "pay-as-you-go"
    releaseType     := r.FormValue("releaseType")
    version         := r.FormValue("version")
    showcase        := strings.Split(r.FormValue("showcase"), ",") // comma-separated values
    tags            := strings.Split(r.FormValue("tags"), ",") // comma-separated values
    runEndpoint     := ""
    pulls           := 0
    license         := r.FormValue("license")
    verified        := false

    if err != nil {
        return models.Tool{}, err
    }

    // Generate error if required fields are empty
    if len(title) == 0 || 
        len(thumbnail) == 0 || 
        len(developer) == 0 || 
        len(pricingType) == 0 || 
        len(releaseType) == 0 || 
        len(version) == 0 || 
        len(showcase) == 0 || 
        len(license) == 0 {

        log.Error("Missing value make sure that all product fields are present") 
        return models.Tool{}, errors.New("Missing value in form")
    }

    return models.Tool{
        UID:          uuid.New().String(),
        Title:        title,
        Description:  description,
        Rating:       rating,
        Price:        price,
        Thumbnail:    thumbnail,
        Overview:     overview,
        CommunityURL: communityURL,
        Reviews:      reviews,
        Developer:    developer,
        PricingType:  pricingType,
        ReleaseType:  releaseType,
        Version:      version,
        Showcase:     showcase,
        Tags:         tags,
        RunEndpoint:  runEndpoint,
        Pulls:        pulls,
        License:      license,
        Verified:     verified,
    }, nil
}

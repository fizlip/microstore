// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
    "mime/multipart"
    "context"
    "fmt"
    "crypto/rand"
    "encoding/base64"
    
    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"

    log "github.com/sirupsen/logrus"
)

// Vendor will handle the business logic related to vendors
type Vendor struct {
    Name                string
	Email               string
	Rating              float64
	ProductsCreated     int
	RequestHandled      int64
	ProductCatalog      []models.Tool
	PaymentExecutor     models.PaymentExecutor
    ToolStorage         models.DatabaseReadWriter[models.Tool]
    AnalyticStorage     models.DatabaseReadWriter[models.ProductAnalytics]
    APIKeyStorage       models.DatabaseReadWriter[models.APIKey]
}

// VendorController is an interface that all Vendors implement
type VendorController interface {
    GetTools() ([]models.Tool, error)
    CreateTool() error
}

func (v *Vendor) Login() {
	fmt.Println("Login: " + v.Email)
}

func (v *Vendor) Register() {
	fmt.Println("Register: " + v.Email)
}

func (v *Vendor) Logout() {
	fmt.Println("Logout: " + v.Email)
}

// GetTools will return an array of all tools created by vendor
func (v *Vendor) GetTools() ([]models.Tool, error){
    return v.ToolStorage.ReadCol("developer", v.Name)
}

// Create will create a tool by inserting it into the data store
func (v *Vendor) CreateTool(tool models.Tool) error {
    return v.ToolStorage.Create(tool, productToAttributeMapConverter)
}

// CreateAPIKey will create an api key by inserting it into the data store
func (v *Vendor) CreateAPIKey(username string) (string, error) {
    var apiKey models.APIKey
    var err error
    apiKey, err  =  models.NewAPIKey(username)
    if err != nil {
        return "", err
    }

    return apiKey.Key, v.APIKeyStorage.Create(apiKey, apiToAttrMap)
}

// DeleteTool will delete a tool
func (v *Vendor) DeleteTool(uti string) error {
    log.Println("Deleting: ", uti)
    return v.ToolStorage.DeleteGSIItem("uti-index", "uti", uti)
}

// DeleteAPIKey will delete an api key 
func (v *Vendor) DeleteAPIKey(username, partialKey string) error {
    log.Println("Deleting key")
    return v.APIKeyStorage.DeleteDetailed(username, partialKey)
}

// ReadAllAPIKeys will read all api keys of a username, we remove the hashed api key
// for security
func (v *Vendor) ReadAllAPIKeys(username string) ([]models.APIKey, error) {

    keys, err := v.APIKeyStorage.Read(username)

    // Remove hashed value when reading users api keys, not needed by user
    for i, _ := range keys {
        keys[i].Key = ""
    }
    return keys, err
}

// ReadAPIKey will read api key of a username
func (v *Vendor) ReadAPIKey(username, key string) (models.APIKey, error) {
    return v.APIKeyStorage.ReadSpecific(username, key)
}

// ReadAPIKey will read api key of a username
func (v *Vendor) ValidateAPIKey(apiKey, storedKey string) bool {
    return models.ValidateKey(apiKey, storedKey)
}

// ExecuteTool ...
func (v *Vendor) ExecuteTool(cid string) {
	log.Println("Give moeny to vendor")
}

// UploadShowcaseImage will upload a file to an S3 bucket 
func (v *Vendor) UploadShowcaseImage(tid string, filename string, replace string, image multipart.File) error {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        return fmt.Errorf("unable to load SDK config: %v", err)
    }

    bucketName := "f4-public"
    objectKey := filename

    // Create an S3 client
    client := s3.NewFromConfig(cfg)
    bucket := NewBucketManager(client, bucketName)
    if replace == "true" {
        bucket.ClearFolder(context.TODO(), "showcases/" + tid)
    }

    err = bucket.PutObject(context.TODO(), "showcases/" + tid + "/" + objectKey, image)

    if err != nil {
        log.Error("Failed to upload image: ", err)
        return err
    } else {
        log.Println("Uploaded: ", tid)
    }

    log.Println("Showcase updated")
	return nil
}

// GetToolUsage will return the analytics for a tool that the vendor has created
func (v *Vendor) GetToolUsage() ([]models.ProductAnalytics, error) {
	vendorTools, err := v.GetTools()
    if err != nil {
        return []models.ProductAnalytics{}, err
    }

	res := []models.ProductAnalytics{}

	for _, product := range vendorTools {
		analytics, err := v.AnalyticStorage.Read(product.UID)
		if err != nil {
			log.Error("Cuould not read", product.UID)
		} else {
			res = append(res, analytics...)
		}
	}
	return res, nil
}

// NewVendorController is a factory method that initialised a data store and
// returns a Vendor struct. If you need to change the data store for some reason
// you only need to update this method.
func NewVendorController(name string) (*Vendor, error){
    // Init tools data store connection
    toolT, err := models.NewDynamoTable[models.Tool]("f4-tools-prod", "uid", "")
    if err != nil {
        return &Vendor{}, err 
    }
    // Init analytics data store connection
    analyticsT, err := models.NewDynamoTable[models.ProductAnalytics]("f4-product-analytics", "product_id", "timestamp")
    if err != nil {
        return &Vendor{}, err 
    }

    apiKeyT, err := models.NewDynamoTable[models.APIKey]("f4-api-keys", "username", "key")
    if err != nil {
        return &Vendor{}, err 
    }
    
    vendor := Vendor{
        Name: name, 
        ToolStorage: toolT, 
        AnalyticStorage: analyticsT, 
        APIKeyStorage: apiKeyT,
    }
    return &vendor, nil
}

func generateSecureString(length int) (string, error) {
    bytes := make([]byte, length)
    
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    
    key := base64.RawURLEncoding.EncodeToString(bytes)

    if len(key) > length {
        key = key[:length]
    }
    
    return key, nil 
} 

func apiToAttrMap(apikey models.APIKey) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "key": &types.AttributeValueMemberS{Value: apikey.HashedKey},
        "username": &types.AttributeValueMemberS{Value: apikey.Username},
        "created": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", apikey.Created)},
        "partial": &types.AttributeValueMemberS{Value: apikey.Key[0:2] + "..." + apikey.Key[39:]},
    }
}

// Utils
func productToAttributeMapConverter(tool models.Tool) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "uid": &types.AttributeValueMemberS{Value: tool.UID},
        "developer": &types.AttributeValueMemberS{Value: tool.Developer},
        "title": &types.AttributeValueMemberS{Value: tool.Title},
        "uti": &types.AttributeValueMemberS{Value: tool.UTI},
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
        "endpoint_names": &types.AttributeValueMemberSS{Value: tool.EndpointNames},
        "endpoint_descriptions": &types.AttributeValueMemberSS{Value: tool.EndpointDescriptions},
        "endpoint_parameters": &types.AttributeValueMemberSS{Value: tool.EndpointParameters},
    }
}

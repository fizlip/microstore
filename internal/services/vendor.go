package services

import (
    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    "fmt"

    log "github.com/sirupsen/logrus"
)

type Vendor struct {
    Name                string
	Email               string
	Rating              float64
	ProductsCreated     int
	RequestHandled      int64
	ProductCatalog      []models.Tool
	PaymentExecutor     models.PaymentExecutor
    ToolStorage         models.DatabaseReadWriter[models.Tool]
    AnalyticStorage    models.DatabaseReadWriter[models.ProductAnalytics]
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
    return v.ToolStorage.Read(v.Name)
}

// Create will create a tool by inserting it into the data store
func (v *Vendor) CreateTool(tool models.Tool) error {
    return v.ToolStorage.Create(tool, productToAttributeMapConverter)
}

func (v *Vendor) ExecuteTool(cid string) {
	log.Println("Give moeny to vendor")
}

func (v *Vendor) GetToolUsage() ([]models.ProductAnalytics, error) {
	vendorTools, err := v.GetTools()
    if err != nil {
        return []models.ProductAnalytics{}, err
    }

	res := []models.ProductAnalytics{}

	for _, product := range vendorTools {
		fmt.Println(product.UID)
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
    toolT, err := models.NewMockDynamoTable[models.Tool]("f4-tools", "developer", "uid")
    if err != nil {
        return &Vendor{}, err 
    }
    analyticsT, err := models.NewMockDynamoTable[models.ProductAnalytics]("f4-tools", "developer", "uid")
    if err != nil {
        return &Vendor{}, err 
    }
    vendor := Vendor{Name: name, ToolStorage: toolT, AnalyticStorage: analyticsT}
    return &vendor, nil
}

// Utils
func productToAttributeMapConverter(tool models.Tool) map[string]types.AttributeValue {
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

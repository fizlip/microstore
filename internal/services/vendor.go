package services

import (
    "github.com/fizlip/microstore/internal/models"
    "fmt"

    log "github.com/sirupsen/logrus"
)

type Vendor struct {
    Name            string
	Email           string
	Rating          float64
	ProductsCreated int
	RequestHandled  int64
	ProductCatalog  []models.Tool
	PaymentExecutor models.PaymentExecutor
    ToolStorage     models.DatabaseReadWriter[models.Tool]
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
    return v.ToolStorage.Create(tool)
}

func (v *Vendor) ExecuteTool(cid string) {
	log.Println("Give moeny to vendor")
}

func (v *Vendor) GetToolUsage() ([]models.ProductAnalytics, error) {
	vendorTools, err := v.GetTools()
    if err != nil {
        return []models.ProductAnalytics{}, err
    }
	pa, err := NewProductAnalyticsTable("product_analytics")
	if err != nil {
		return []models.ProductAnalytics{}, err
	}

	res := []models.ProductAnalytics{}

	for _, product := range vendorTools {
		fmt.Println(product.UID)
		analytics, err := (*pa).Read(product.UID)
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
func NewVendorController(name string) (Vendor, error){
    table, err := models.CreateToolsTable()
    if err != nil {
        return Vendor{}, err 
    }
    vendor      := Vendor{Name: name, ToolStorage: table}
    return vendor, nil
}

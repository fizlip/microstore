package services

import (
    "github.com/fizlip/microstore/internal/models"
)


// VendorController is an interface that all Vendors implement
type VendorController interface {
    GetTools() ([]models.Tool, error)
    CreateTool() error
}

// Vendor represents a user that sells a product on f4rmhouse. It consits
// of a username and a data store used to interact with the tools sold 
// by the vendor 
type Vendor struct {
    name string
    toolStorage models.DatabaseReadWriter[models.Tool]
}

// GetTools will return an array of all tools created by vendor
func (v *Vendor) GetTools() ([]models.Tool, error){
    return v.toolStorage.Read(v.name)
}

// Create will create a tool by inserting it into the data store
func (v *Vendor) CreateTool(tool models.Tool) error {
    return v.toolStorage.Create(tool)
}

// NewVendorController is a factory method that initialised a data store and
// returns a Vendor struct. If you need to change the data store for some reason
// you only need to update this method.
func NewVendorController(name string) (Vendor, error){
    table, err := models.CreateToolsTable()
    if err != nil {
        return Vendor{}, err 
    }
    vendor      := Vendor{name: name, toolStorage: table}
    return vendor, nil
}

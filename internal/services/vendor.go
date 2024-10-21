package services

import (
    "github.com/fizlip/microstore/internal/models"
)

type VendorController interface {
    GetTools() ([]models.Tool, error)
    CreateTool() error
}

type Vendor struct {
    name string
    toolStorage models.DatabaseReadWriter[models.Tool]
}

func (v *Vendor) GetTools() ([]models.Tool, error){
    return v.toolStorage.Read(v.name)
}

func (v *Vendor) CreateTool(tool models.Tool) error {
    return v.toolStorage.Create(tool)
}

func NewVendorController(name string) (Vendor, error){
    table, err := models.CreateToolsTable()
    if err != nil {
        return Vendor{}, err 
    }
    vendor      := Vendor{name: name, toolStorage: table}
    return vendor, nil
}

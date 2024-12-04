package services

import (
    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    // log "github.com/sirupsen/logrus"
)

type ToolNameService struct {
    NameServiceStorage  models.DatabaseReadWriter[models.ToolName]
}

func (t *ToolNameService) GetTool(tool string) (models.ToolName, error){
    return t.NameServiceStorage.ReadSpecific(tool, "")
}

func (t *ToolNameService) CreateRecord(tool models.ToolName) error {
    return t.NameServiceStorage.Create(tool, toolToAttributeMapConverter)
}

func (t *ToolNameService) DeleteRecord(uti string) error {
    return t.NameServiceStorage.Delete(uti)
}

func (t *ToolNameService) InsertIP(uti string) error {
    return t.NameServiceStorage.UpdateAll(uti, "SS,ips", "123,456")
}

// NewVendorController is a factory method that initialised a data store and
// returns a Vendor struct. If you need to change the data store for some reason
// you only need to update this method.
func NewToolNameService() (*ToolNameService, error){
    // Init tools data store connection
    toolNameT, err := models.NewDynamoTable[models.ToolName]("f4-ns", "uti", "")
    if err != nil {
        return &ToolNameService{}, err 
    }
    
    ns := ToolNameService{
        NameServiceStorage: toolNameT, 
    }
    return &ns, nil
}

// Utils
func toolToAttributeMapConverter(tool models.ToolName) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "uti": &types.AttributeValueMemberS{Value: tool.UTI},
    }
}

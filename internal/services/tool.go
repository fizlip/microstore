// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
	// log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fizlip/microstore/internal/models"
)

type ToolNameService struct {
	NameServiceStorage models.DatabaseReadWriter[models.ToolName]
}

func (t *ToolNameService) GetTool(tool string) ([]models.ToolName, error) {
	return t.NameServiceStorage.Read(tool)
}

func (t *ToolNameService) CreateRecord(tool models.ToolName) error {
	return t.NameServiceStorage.Create(tool, toolToAttributeMapConverter)
}

func (t *ToolNameService) DeleteRecord(uti string) error {
	return t.NameServiceStorage.Delete(uti)
}

func (t *ToolNameService) InsertIP(uti, ips string) error {
	return t.NameServiceStorage.UpdateAll(uti, "SS,ips", ips)
}

// NewVendorController is a factory method that initialised a data store and
// returns a Vendor struct. If you need to change the data store for some reason
// you only need to update this method.
func NewToolNameService() (*ToolNameService, error) {
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
		"uti":             &types.AttributeValueMemberS{Value: tool.UTI},
		"endpoints":       &types.AttributeValueMemberSS{Value: tool.Endpoints},
		"descriptions":    &types.AttributeValueMemberSS{Value: tool.Descriptions},
		"parameters":      &types.AttributeValueMemberSS{Value: tool.Parameters},
		"deployment_type": &types.AttributeValueMemberS{Value: tool.DeploymentType},
	}
}

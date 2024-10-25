package services

import (

    "time"    

    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    log "github.com/sirupsen/logrus"
)

type Customer struct {
    Name                string
    analyticsStore      models.DatabaseReadWriter[models.UserAnalytics]
    toolAnalyticsStore  models.DatabaseReadWriter[models.ProductAnalytics]
}

// GetAnalytics returns all product analytics data related to the customer for 
// the current month
func (c *Customer) GetAnalytics() ([]models.UserAnalytics, error) {
    return c.analyticsStore.Read(c.Name) 
}

// GetUsage returns all products that user has used for the past month
// and how many times they have been used.
func (c *Customer) GetUsage() ([]models.ProductAnalytics, error) {
    return []models.ProductAnalytics{}, nil
}

func (c *Customer) RunProduct(uid string, tid string, data map[string]interface{}) (string, error) {
    evt := models.ProductAnalytics{
		ProductID:    tid,
        Status: "400",
		Consumer: uid,
    }
    log.Println("Execute tool with data = ", data)
    log.Println("Run tool sucessfully")
    evt.Status = "200"
    evt.Timestamp = time.Now().Format("2006-01-02 15:04:05")
    c.toolAnalyticsStore.Create(evt, evtToAttributeMapConverter)
    log.Println("Create payment event")
    
    return "Hello world!", nil
}

// NewCustomer returns a customer struct with the given username.
// the datastore can be changed here.
func NewCustomer(username string) (*Customer, error) {
    analyticsT, err := models.NewMockDynamoTable[models.UserAnalytics]("f4-product-analytics", "developer", "uid")
    if err != nil {
        return &Customer{}, err 
    }
    toolAnalyticsT, err := models.NewMockDynamoTable[models.ProductAnalytics]("f4-product-analytics", "developer", "uid")
    if err != nil {
        return &Customer{}, err 
    }
    return &Customer{
        Name: username, 
        analyticsStore: analyticsT, 
        toolAnalyticsStore: toolAnalyticsT,
    }, nil
}

// Utils
func evtToAttributeMapConverter(p models.ProductAnalytics) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "pid": &types.AttributeValueMemberS{Value: p.ProductID},
        "timestamp": &types.AttributeValueMemberS{Value: p.Timestamp},
        "status": &types.AttributeValueMemberS{Value: p.Status},
        "consumer": &types.AttributeValueMemberS{Value: p.Consumer},
    }
}

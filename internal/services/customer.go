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
    toolStore           models.DatabaseReadWriter[models.Tool]
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

func (c *Customer) incrementUsage(uid, tid, time string) error {
    // Gives the current days analytics
    analytics, err := c.analyticsStore.ReadSpecific(uid, time)
    if err != nil {
        return err
    }
    if analytics.CustomerID == "" {
        analytics = models.UserAnalytics{
            CustomerID: uid,
            Timestamp: time,
            PaidRequests: 0,
            FreeRequests: 0,
            Products: map[string]models.ProductData{}, 
        }
    }

    // How do I make sure only one item returne? Need to add check when creating so that
    // all ids must be unique
    tools, err := c.toolStore.Read(tid)
    if err != nil {
        return err
    }
    if tools[0].Price == 0.0 {
        analytics.FreeRequests += 1
    } else {
        analytics.PaidRequests += 1
    }
    
    p, ok := analytics.Products[tid]
    if ok {
       p.NumRequests += 1 
    } else {
        analytics.Products[tid] = models.ProductData{
            ProductID: tid, 
            ProductName: tools[0].Title, 
            Price: tools[0].Price,
            NumRequests: 1,
        }
    }
    
    c.analyticsStore.Create(analytics, userAnalyticsToAttributeMapConverter)
    return nil
    
}

func (c *Customer) RunProduct(uid string, tid string, data map[string]interface{}) (string, error) {
    evt := models.ProductAnalytics{
		ProductID: tid,
        Status: "400",
		Consumer: uid,
    }
    log.Println("Execute tool with data = ", data)
    log.Println("Run tool sucessfully")
    before, _ := c.analyticsStore.Read("2024-10-25")
    log.Println("Analytics before run ", before)
    evt.Status = "200"
    evt.Timestamp = time.Now().Format("2006-01-02 15:04:05")
    c.toolAnalyticsStore.Create(evt, evtToAttributeMapConverter)

    c.incrementUsage(uid, tid, time.Now().Format("2006-01-02"))
    after, _ := c.analyticsStore.Read("2024-10-25")
    log.Println("Analytics before run ", after)

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

    toolT, err := models.NewMockDynamoTable[models.Tool]("f4-tools", "developer", "uid")
    if err != nil {
        return &Customer{}, err 
    }

    return &Customer{
        Name: username, 
        analyticsStore: analyticsT, 
        toolAnalyticsStore: toolAnalyticsT,
        toolStore: toolT,
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

func userAnalyticsToAttributeMapConverter(u models.UserAnalytics) map[string]types.AttributeValue {
    return map[string]types.AttributeValue{
        "customer_id": &types.AttributeValueMemberS{Value: u.CustomerID},
        "timestamp": &types.AttributeValueMemberS{Value: u.Timestamp},
    }
}

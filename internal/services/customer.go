package services

import (

    "time"    
    "errors"

    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

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
        analytics.Products[tid] = p
    } else {
        analytics.Products[tid] = models.ProductData{
            ProductID: tid, 
            ProductName: tools[0].Title, 
            Price: tools[0].Price,
            NumRequests: 1,
        }
    }
    
    err = c.analyticsStore.Create(analytics, userAnalyticsToAttributeMapConverter)
    if err != nil {
        return err
    }

    return nil
    
}

func (c *Customer) RunProduct(uid string, tid string, data map[string]interface{}) (string, error) {
    var err error
    tool, err := c.toolStore.Read(tid)
    log.Println("tid", tid)
    if err != nil {
        return "", errors.New("The tool you want to run does not exist.")
    } 
    log.Println("tool", tool[0])
    evt := models.ProductAnalytics{
		ProductID: tid,
        Status: "400",
		Consumer: uid,
    }

    evt.Status = "200"
    evt.Timestamp = time.Now().Format("2006-01-02 15:04:05")
    err = c.toolAnalyticsStore.Create(evt, evtToAttributeMapConverter)
    if err != nil {
        return "", err
    }

    err = c.incrementUsage(uid, tid, time.Now().Format("2006-01-02"))
    if err != nil {
        return "", err
    }

    return "Hello world!", nil
}

// NewCustomer returns a customer struct with the given username.
// the datastore can be changed here.
func NewCustomer(username string) (*Customer, error) {
    analyticsT, err := models.NewDynamoTable[models.UserAnalytics]("f4-user-analytics", "customer_id", "timestamp")
    if err != nil {
        return &Customer{}, err 
    }

    toolAnalyticsT, err := models.NewDynamoTable[models.ProductAnalytics]("f4-product-analytics", "product_id", "timestamp")
    if err != nil {
        return &Customer{}, err 
    }

    toolT, err := models.NewDynamoTable[models.Tool]("f4-tools-prod", "uid", "")
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
        "product_id": &types.AttributeValueMemberS{Value: p.ProductID},
        "timestamp": &types.AttributeValueMemberS{Value: p.Timestamp},
        "status": &types.AttributeValueMemberS{Value: p.Status},
        "consumer": &types.AttributeValueMemberS{Value: p.Consumer},
    }
}

func userAnalyticsToAttributeMapConverter(u models.UserAnalytics) map[string]types.AttributeValue {
    av, _ := attributevalue.MarshalMap(u)
    return av
    //return map[string]types.AttributeValue{
    //    "customer_id": &types.AttributeValueMemberS{Value: u.CustomerID},
    //    "timestamp": &types.AttributeValueMemberS{Value: u.Timestamp},
    //    "free_requests": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", u.FreeRequests)},
    //    "paid_requests": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", u.PaidRequests)},
    //    "productsl": &types.AttributeValueMemberSS{Value: fmt.Sprintf("%d", u.PaidRequests)},

    //}
}

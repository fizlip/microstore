// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (

    "time"    
    "strconv"
    "fmt"
    "errors"
    "net/http"
    "io/ioutil"
    "bytes"

    "github.com/fizlip/microstore/internal/models"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

    log "github.com/sirupsen/logrus"
)

// Customer
type Customer struct {
    Name                string
    analyticsStore      models.DatabaseReadWriter[models.UserAnalytics]
    toolAnalyticsStore  models.DatabaseReadWriter[models.ProductAnalytics]
    toolStore           models.DatabaseReadWriter[models.Tool]
    f4rmerStore         models.DatabaseReadWriter[models.F4rmer]
}

// Run ...
func (c *Customer) Run(uti, endpoint, params string) (string, error) {
    url := fmt.Sprintf("http://%s:8080/%s", uti, endpoint)

    data := []byte(params)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    // resp, err := http.Post("http://" + uti + ":8080/" + endpoint + "?" + params)
    if err != nil {
        log.Println("Error making GET request: %v\n", err)
        return "", err
    }
    // Set headers
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return "", err
    }

    defer resp.Body.Close()

    // Read response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response: %v\n", err)
        return "", err
    }

    log.Println("Action has been executed with status 200")

    return string(body), nil
}

func (c *Customer) CreateEmptyAnalyticsRecord() error {
    emptyObject := models.UserAnalytics{}
    emptyObject.CustomerID = c.Name
    emptyObject.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
    return c.analyticsStore.Create(emptyObject, userAnalyticsToAttributeMapConverter)
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

// DeleteF4rmer will delete a f4rmer
func (c *Customer) DeleteAnalytics() (error) {
    analytics, err := c.analyticsStore.Read(c.Name)
    if err != nil {
        log.Error("Could not read analytics for ", c.Name)
        return err
    }
    if len(analytics) == 0 {
        log.Error("There are no records with the name: ", c.Name)
        return errors.New("Nothing to delete!") 
    }
    return c.analyticsStore.DeleteDetailed(analytics[0].CustomerID, analytics[0].Timestamp)
}

// incrementUsage is a private function that updates the user analytics of
// a customer that uses a tool.
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

// RunProduct will execute a tool and also update the user analytics
func (c *Customer) RunProduct(uid string, tid string, data map[string]interface{}) (string, error) {
    var err error
    tool, err := c.toolStore.Read(tid)
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

// CreateF4rmer will create a new f4rmer
func (c *Customer) CreateF4rmer(f4rmer models.F4rmer) (error) {
    return c.f4rmerStore.Create(f4rmer, f4rmerToAttributeMapConverter)
}

// ReadF4rmers will read all f4rmers that have been created by a customer
func (c *Customer) ReadF4rmers() ([]models.F4rmer, error) {
    return c.f4rmerStore.Read(c.Name)
}

// ReadF4rmer will read one single f4rmer
func (c *Customer) ReadF4rmer(username, uid string) (models.F4rmer, error) {
    return c.f4rmerStore.ReadSpecific(username, uid)
}

// UpdateF4rmer will update a f4rmer
func (c *Customer) UpdateF4rmer(creator, uid, col, val string) (error) {
    return c.f4rmerStore.Update(creator, uid, col, val)
}

// DeleteF4rmer will delete a f4rmer
func (c *Customer) DeleteF4rmer(creator, uid string) (error) {
    return c.f4rmerStore.DeleteDetailed(creator, uid)
}

// AddTool will add a tool to a f4rmer
func (c *Customer) AddTool(username, uid, tid string) (error) {

    f4rmer, err := c.ReadF4rmer("fillatino@gmail.com", uid)
    if err != nil {
        log.Println("Cannot read f4rmer")
        return err
    }

    // Make sure tool exists
    tool, err := c.toolStore.Read(tid) 
    if err != nil {
        log.Println("Cannot read tool")
        return err
    }
    if contains(f4rmer.Toolbox, tool[0]) {
        log.Println("Tool already exists in toolbox. Tool is not added.")
        return errors.New("Tool already in toolbox") 
    }

    newToolbox := append(f4rmer.Toolbox, tool[0]) 
    f4rmer.Toolbox = newToolbox
    
    err = c.CreateF4rmer(f4rmer)
    if err != nil {
        log.Println("Cannot add tool")
        return err
    }

    return nil
}

// constains is a utility function to check it a slice of tools contains a specific tool
func contains(ts []models.Tool, t models.Tool) bool {
    for _, item := range ts {
        if item.UID == t.UID {
            return true
        }
    }
    return false
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
    f4rmerT, err := models.NewDynamoTable[models.F4rmer]("f4-f4rmers", "creator", "uid")
    if err != nil {
        return &Customer{}, err 
    }

    return &Customer{
        Name: username, 
        analyticsStore: analyticsT, 
        toolAnalyticsStore: toolAnalyticsT,
        toolStore: toolT,
        f4rmerStore: f4rmerT,
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

func f4rmerToAttributeMapConverter(f4 models.F4rmer) map[string]types.AttributeValue {
    av, _ := attributevalue.MarshalMap(f4)
    return av
}

func userAnalyticsToAttributeMapConverter(u models.UserAnalytics) map[string]types.AttributeValue {
    av, _ := attributevalue.MarshalMap(u)
    return av
}

package models

import (
    "fmt"
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    log "github.com/sirupsen/logrus"
)

type ToolsTable struct{
    svc *dynamodb.Client
}

func (t *ToolsTable) Read(id string) ([]Tool, error) {
    return []Tool{}, nil
}

func (t *ToolsTable) ReadSpecific(pk string, sk string) (Tool, error) {
    return Tool{}, nil
}

func (t *ToolsTable) Update(pk, sk string) error {
    return nil
}

// Create will insert an item into the tools database
func (t *ToolsTable) Create(tool Tool) error {
    fmt.Println("Creating: ", tool)

    // Is this there a better wat to do this??
    item := map[string]types.AttributeValue{
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
    }

    input := &dynamodb.PutItemInput{
        TableName: aws.String("f4-tools"),
        Item:      item,
    }

    // Insert the item
    _, err := t.svc.PutItem(context.TODO(), input)
    if err != nil {
        log.Fatalf("failed to put item in DynamoDB, %v", err)
        return err
    }

    return nil
}

// Delete shold be used rarely even though pk should be unique in some
// db implementations it  can delete more than one time
func (t *ToolsTable) Delete(pk string) (error) {
    return nil
}

// DelteDetailed will delete the item with the given pk and sk. This method should be
// prefered over the simpler Delete method
func (t *ToolsTable) DeleteDetailed(pk, sk string) error {
    return nil
}

// CreateToolsTable will return an implementation of a table for the Tools dynamoDB  
// table that is used to interact with tools in the f4rmhouse ecosystem.
func CreateToolsTable() (*ToolsTable, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    var table ToolsTable = ToolsTable{}

    if err != nil {
        log.Error("unable to load SDK config, %v", err)
        return &table, err
    }

    table = ToolsTable{svc: dynamodb.NewFromConfig(cfg)}
    return &table, nil
}

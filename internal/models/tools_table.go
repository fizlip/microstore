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

type Product struct {
	UID          string   `json:"uid"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Rating       float64  `json:"rating"`
	Price        float64  `json:"price"`
	Thumbnail    string   `json:"thumbnail"`
	Overview     string   `json:"overview"`
	CommunityURL string   `json:"communityURL"`
	Reviews      int      `json:"reviews"`
	Developer    string   `json:"developer"`
	PricingType  string   `json:"pricingType"`
	ReleaseType  string   `json:"releaseType"`
	Version      string   `json:"version"`
	Showcase     []string `json:"showcase"`
	Tags         []string `json:"tags"`
	RunEndpoint  string   `json:"run_endpoint"`
	Pulls        int      `json:"pulls"`
    License      string   `json:"license"`
}

type ToolsTable struct{
    svc *dynamodb.Client
}

func (t *ToolsTable) Read(id string) ([]Product, error) {
    return []Product{}, nil
}

func (t *ToolsTable) ReadSpecific(pk string, sk string) (Product, error) {
    return Product{}, nil
}

func (t *ToolsTable) Update(pk, sk string) error {
    return nil
}

func (t *ToolsTable) Create(p Product) error {
    fmt.Println("Creating: ", p)
    item := map[string]types.AttributeValue{
        "uid": &types.AttributeValueMemberS{Value: p.UID},
        "developer": &types.AttributeValueMemberS{Value: p.Developer},
        "title": &types.AttributeValueMemberS{Value: p.Title},
        "description": &types.AttributeValueMemberS{Value: p.Description},
        "rating": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", p.Rating)},
        "Pricing": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", p.Price)},
        "thumbnail": &types.AttributeValueMemberS{Value: p.Thumbnail},
        "overview": &types.AttributeValueMemberS{Value: p.Overview},
        "communityURL": &types.AttributeValueMemberS{Value: p.CommunityURL},
        "reviews": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", p.Reviews)},
        "pricingType": &types.AttributeValueMemberS{Value: p.PricingType},
        "releaseType": &types.AttributeValueMemberS{Value: p.ReleaseType},
        "version": &types.AttributeValueMemberS{Value: p.Version},
        "showcase": &types.AttributeValueMemberSS{Value: p.Showcase},
        "tags": &types.AttributeValueMemberSS{Value: p.Tags},
        "runEndpoint": &types.AttributeValueMemberS{Value: p.RunEndpoint},
        "pulls": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", p.Pulls)},
        "license": &types.AttributeValueMemberS{Value: p.License},
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

func (t *ToolsTable) Delete(pk string) (error) {
    return nil
}

func (t *ToolsTable) DeleteDetailed(pk, sk string) error {
    return nil
}

func CreateToolsTable() (*Table[Product], error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    var table Table[Product] = &ToolsTable{}

    if err != nil {
        log.Error("unable to load SDK config, %v", err)
        return &table, err
    }

    table = &ToolsTable{svc: dynamodb.NewFromConfig(cfg)}
    return &table, nil
}

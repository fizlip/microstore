package models

import (
    "fmt"
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

    log "github.com/sirupsen/logrus"
)

type DynamoTable[T any] struct{
    pk          string
    sk          string
    tableName   string
    svc         *dynamodb.Client
}

func (d *DynamoTable[T]) Read(id string) ([]T, error) {
    input := &dynamodb.QueryInput{
        TableName: aws.String(d.tableName),
        KeyConditionExpression: aws.String(fmt.Sprintf("%s = :v", d.pk)),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":v": &types.AttributeValueMemberS{Value: id},
        },
    }

    // Call the Query API
    result, err := d.svc.Query(context.TODO(), input)
    if err != nil {
        log.Fatalf("failed to query id, %v", err)
        return []T{}, err
    }

    var ts []T
    err = attributevalue.UnmarshalListOfMaps(result.Items, &ts)
    if err != nil {
        log.Fatalf("Failed to unmarshal Query result into []Tool: %v", err)
        return []T{}, err
    }


    return ts, nil
}

func (d *DynamoTable[T]) ReadSpecific(pk string, sk string) (T, error) {
    var t T
    return t, nil
}

func (d *DynamoTable[T]) Update(pk, sk string) error {
    return nil
}

// Create will insert an item into the tools database
func (d *DynamoTable[T]) Create(t T, f func (t T) map[string]types.AttributeValue) error {
    // Convert T into format understandable by &dynamodb
    item := f(t)

    input := &dynamodb.PutItemInput{
        TableName: aws.String(d.tableName),
        Item:      item,
    }

    // Insert the item
    _, err := d.svc.PutItem(context.TODO(), input)
    if err != nil {
        return err
    }

    return nil
}

// Delete shold be used rarely even though pk should be unique in some
// db implementations it  can delete more than one time
func (d *DynamoTable[T]) Delete(pk string) (error) {
    return nil
}

// DelteDetailed will delete the item with the given pk and sk. This method should be
// prefered over the simpler Delete method
func (d *DynamoTable[T]) DeleteDetailed(pk, sk string) error {
    key := map[string]types.AttributeValue{
        d.pk: &types.AttributeValueMemberS{Value: pk},
        d.sk: &types.AttributeValueMemberS{Value: sk},
    }
    log.Println(pk)
    log.Println(sk)
    log.Println(d.tableName)
    log.Println(key)
    input := &dynamodb.DeleteItemInput{
        TableName: aws.String("f4-api-keys"),
        Key:       key,
    }

     res, err := d.svc.DeleteItem(context.TODO(), input)
    if err != nil {
        log.Error("Could not delete item from table")
        return err
    }
    log.Println(res)
    return nil
}

// CreateToolsTable will return an implementation of a table for the Tools dynamoDB  
// table that is used to interact with tools in the f4rmhouse ecosystem.
func NewDynamoTable[T any](tableName, pk, sk string) (*DynamoTable[T], error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    var table DynamoTable[T] = DynamoTable[T]{pk: pk, sk: sk, tableName: tableName}

    if err != nil {
        log.Error("unable to load SDK config, %v", err)
        return &table, err
    }

    table.svc = dynamodb.NewFromConfig(cfg)
    return &table, nil
}

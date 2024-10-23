package models

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    log "github.com/sirupsen/logrus"
)

type MockDynamoTable[T any] struct{
    pk          string
    sk          string
    tableName   string
    svc         *dynamodb.Client
}

func (m *MockDynamoTable[T]) Read(id string) ([]T, error) {
    var ts []T
    return ts, nil
}

func (m *MockDynamoTable[T]) ReadSpecific(pk string, sk string) (T, error) {
    var t T
    return t, nil
}

func (m *MockDynamoTable[T]) Update(pk, sk string) error {
    return nil
}

// Create will insert an item into the tools database
func (m *MockDynamoTable[T]) Create(t T, f func (t T) map[string]types.AttributeValue) error {
    return nil
}

// Delete shold be used rarely even though pk should be unique in some
// db implementations it  can delete more than one time
func (m *MockDynamoTable[T]) Delete(pk string) (error) {
    return nil
}

// DelteDetailed will delete the item with the given pk and sk. This method should be
// prefered over the simpler Delete method
func (m *MockDynamoTable[T]) DeleteDetailed(pk, sk string) error {
    return nil
}

// CreateToolsTable will return an implementation of a table for the Tools dynamoDB  
// table that is used to interact with tools in the f4rmhouse ecosystem.
func NewMockDynamoTable[T any](tableName, pk, sk string) (*MockDynamoTable[T], error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    var table MockDynamoTable[T] = MockDynamoTable[T]{pk: pk, sk: sk, tableName: tableName}

    if err != nil {
        log.Error("unable to load SDK config, %v", err)
        return &table, err
    }

    table.svc = dynamodb.NewFromConfig(cfg)
    return &table, nil
}

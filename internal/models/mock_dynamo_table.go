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
    mock        []T
}

func (d *MockDynamoTable[T]) BatchRead(pks, sks []string) ([]T, error) {
    return []T{}, nil
}

func (m *MockDynamoTable[T]) Read(id string) ([]T, error) {
    return m.mock, nil
}

func (m *MockDynamoTable[T]) ReadCol(col, val string) ([]T, error) {
    return m.mock, nil 
}

func (m *MockDynamoTable[T]) ReadSpecific(pk string, sk string) (T, error) {
    var t T
    return t, nil
}

func (m *MockDynamoTable[T]) Update(pk, sk, col, val string) error {
    return nil
}

func (m *MockDynamoTable[T]) UpdateAll(pk, col, val string) error {
    return nil
}

func (m *MockDynamoTable[T]) UpdateObject(pk, sk, col string, t T) error {
    return nil
}

// Create will insert an item into the tools database
func (m *MockDynamoTable[T]) Create(t T, f func (t T) map[string]types.AttributeValue) error {
    m.mock = append(m.mock, t)
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

// ProduceMockData will insert fake data into the mock datastore. We do this
// using a type switch. This function is only used in dev so we don't need to 
// worry that much about performance issues due to type switches.
func (m *MockDynamoTable[T]) ProduceMockData() {
    var zero T
    switch any(zero).(type) {
        case UserAnalytics:
            m.mock = any(MockUserAnalytics).([]T)
        case ProductAnalytics:
            m.mock = any(MockProductAnalytics).([]T)
        case Tool:
            m.mock = any(MockProductCatalog).([]T)
        case Review:
            m.mock = any(MockReviews).([]T)
        case F4rmer:
            m.mock = any(MockF4rmers).([]T)
        default:
            m.mock = []T{}
    }
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
    table.ProduceMockData()

    return &table, nil
}

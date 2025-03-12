// package models contains definitions of common types and models used accross the API.
package models

import (
    "fmt"
    "context"
    "errors"
    "strings"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

    log "github.com/sirupsen/logrus"
)

// DynamoTable is a generic type used to represent a dynamodb table
// All data stores are specific types of this type.
type DynamoTable[T any] struct{
    pk          string
    sk          string
    tableName   string
    svc         *dynamodb.Client
}

// BatchReadPks returns all items that have the given parition keys
func (d *DynamoTable[T]) BatchReadPks(pks []string) ([]T, error) {
    var returnedItems []T
    var keys []map[string]types.AttributeValue

    for _, pk := range pks {
        key := map[string]types.AttributeValue{
            d.pk: &types.AttributeValueMemberS{Value: pk},
        }
        keys = append(keys, key)
    } 
    
    input := &dynamodb.BatchGetItemInput{
        RequestItems: map[string]types.KeysAndAttributes{
            d.tableName: {
                Keys: keys,
            },
        },
    }

    result, err := d.svc.BatchGetItem(context.TODO(), input)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("failed to batch get items: %v", err))
    }

    for _, responses := range result.Responses {
        for _, itemMap := range responses {
            var item T 
            err := attributevalue.UnmarshalMap(itemMap, &item)
            if err != nil {
                return []T{}, errors.New(fmt.Sprintf("failed to unmarshal item: %v", err))
            }
            returnedItems = append(returnedItems, item)
        }
    }
    return returnedItems, nil
}

// BatchRead will return all records with the given partition keys and the given sort keys.
// Use this function when you need fetch more than one item from a table rather than
// making multiple requests.
func (d *DynamoTable[T]) BatchRead(pks, sks []string) ([]T, error) {
    var returnedItems []T
    var keys []map[string]types.AttributeValue
    if len(pks) != len(sks) {
        return []T{}, errors.New("length of partition keys do not match the length of sort keys")
    }
    for index, pk := range pks {
        key := map[string]types.AttributeValue{
            d.pk: &types.AttributeValueMemberS{Value: pk},
            d.sk: &types.AttributeValueMemberS{Value: sks[index]},
        }
        keys = append(keys, key)
    } 
    
    input := &dynamodb.BatchGetItemInput{
        RequestItems: map[string]types.KeysAndAttributes{
            d.tableName: {
                Keys: keys,
            },
        },
    }

    result, err := d.svc.BatchGetItem(context.TODO(), input)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("failed to batch get items: %v", err))
    }

    for _, responses := range result.Responses {
        for _, itemMap := range responses {
            var item T 
            err := attributevalue.UnmarshalMap(itemMap, &item)
            if err != nil {
                return []T{}, errors.New(fmt.Sprintf("failed to unmarshal item: %v", err))
            }
            returnedItems = append(returnedItems, item)
        }
    }
    return returnedItems, nil
}

// ReadCol returns all records where the column (field) is the given value
// i.e. we create the following dynamoDB query: #col = :value
func (d *DynamoTable[T]) ReadCol(col, val string) ([]T, error) {
    var ts []T
    input := &dynamodb.ScanInput{
        TableName: aws.String(d.tableName),
        // Optional: Filter results
        FilterExpression: aws.String("#col = :value"),
        ExpressionAttributeNames: map[string]string{
            "#col": col,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":value": &types.AttributeValueMemberS{Value: val},
        },
    }
    // Call the Query API
    result, err := d.svc.Scan(context.TODO(), input)
    if err != nil {
        log.Fatalf("failed to query id, %v", err)
        return []T{}, err
    }

    err = attributevalue.UnmarshalListOfMaps(result.Items, &ts)
    if err != nil {
        log.Fatalf("Failed to unmarshal Query result into []Tool: %v", err)
        return []T{}, err
    }
    return ts, nil
}

// Read will return an item by the pk. It returns an error if the item does not exists
func (d *DynamoTable[T]) Read(id string) ([]T, error) {
    if d.sk == "" {
        var t T
        key := map[string]types.AttributeValue{
            d.pk: &types.AttributeValueMemberS{Value: id},
        }

        input := &dynamodb.GetItemInput{
            TableName: aws.String(d.tableName),
            Key:       key,
        }
        result, err := d.svc.GetItem(context.TODO(), input)
        if err != nil {
            return []T{}, err
        }
        
        if result.Item == nil {
            return []T{}, errors.New("The requested item does not exist")
        }

        err = attributevalue.UnmarshalMap(result.Item, &t)
        if err != nil {
            return []T{}, err
        }

        return []T{t}, nil
    } else {
        var ts []T
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

        err = attributevalue.UnmarshalListOfMaps(result.Items, &ts)
        if err != nil {
            log.Fatalf("Failed to unmarshal Query result into []Tool: %v", err)
            return []T{}, err
        }
        return ts, nil
    }
    return []T{}, errors.New("Something went wrong reading the value")
}
// ReadAll will return all elements in a table
//
// OBS: only use this when you know this is a small table with <1000 records (to avoid costs) 
func (d *DynamoTable[T]) ReadAll() ([]T, error) {
    var ts []T

    input := &dynamodb.ScanInput{
        TableName: aws.String(d.tableName),
    }

    result, err := d.svc.Scan(context.TODO(), input)
    if err != nil {
        return []T{}, err
    }
    err = attributevalue.UnmarshalListOfMaps(result.Items, &ts)
    if err != nil {
        log.Fatalf("Failed to unmarshal Query result into []Tool: %v", err)
        return []T{}, err
    }

    return ts, nil
}

// ReadSpecific returns a single record with given partition key (pk) and sort key (sk)
func (d *DynamoTable[T]) ReadSpecific(pk string, sk string) (T, error) {
    var t T

    key := map[string]types.AttributeValue{
        d.pk: &types.AttributeValueMemberS{Value: pk},    
        d.sk: &types.AttributeValueMemberS{Value: sk}, 
    }

    input := &dynamodb.GetItemInput{
        TableName: aws.String(d.tableName), // Replace with your table name
        Key:       key,
    }

    result, err := d.svc.GetItem(context.TODO(), input)
    if err != nil {
        return t, err
    }

    err = attributevalue.UnmarshalMap(result.Item, &t)
    if err != nil {
        return t, err
    }

    return t, nil
}

// Update will update the record with the given partition key and sort key
func (d *DynamoTable[T]) Update(pk, sk, col, val string) error {
    key := map[string]types.AttributeValue{
        d.pk: &types.AttributeValueMemberS{Value: pk},    
        d.sk: &types.AttributeValueMemberS{Value: sk}, 
    }
    
    typeColSlice := strings.Split(col, ",")
    if len(typeColSlice) != 2 {
        return errors.New("Column is not formatted correctly. For dynamodb to update you also need to add a type to the column name e.g. S,rating")
    }

    updateType := typeColSlice[0]
    columnName := typeColSlice[1]

    var expressionValue types.AttributeValue
    switch updateType {
        case "S":
            expressionValue = &types.AttributeValueMemberS{Value: val}
        case "N":
            expressionValue = &types.AttributeValueMemberN{Value: val}
    }

    input := &dynamodb.UpdateItemInput{
        TableName: aws.String(d.tableName),
        Key: key,
        UpdateExpression: aws.String("SET #attr = :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": columnName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": expressionValue,
        },
    }
    
    _, err := d.svc.UpdateItem(context.TODO(), input)
    if err != nil {
        return err
    }

    return nil
}

// UpdateAll will update all records with the given private key
func (d *DynamoTable[T]) UpdateAll(pk, col, val string) error {
    key := map[string]types.AttributeValue{
        d.pk: &types.AttributeValueMemberS{Value: pk},    
    }
    
    typeColSlice := strings.Split(col, ",")
    if len(typeColSlice) != 2 {
        return errors.New("Column is not formatted correctly. For dynamodb to update you also need to add a type to the column name e.g. S,rating")
    }

    updateType := typeColSlice[0]
    columnName := typeColSlice[1]

    var expressionValue types.AttributeValue
    switch updateType {
        case "S":
            expressionValue = &types.AttributeValueMemberS{Value: val}
        case "SS":
            expressionValue = &types.AttributeValueMemberSS{Value: strings.Split(val, ",")}
        case "N":
            expressionValue = &types.AttributeValueMemberN{Value: val}
    }

    input := &dynamodb.UpdateItemInput{
        TableName: aws.String(d.tableName),
        Key: key,
        UpdateExpression: aws.String("SET #attr = :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": columnName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": expressionValue,
        },
    }
    
    _, err := d.svc.UpdateItem(context.TODO(), input)
    if err != nil {
        return err
    }

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
    _, err := d.svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
        TableName: aws.String(d.tableName),
        Key: map[string]types.AttributeValue{
            d.pk: &types.AttributeValueMemberS{Value: pk},
        },
    })

    if err != nil {
        log.Error("Could not delete item from table")
        return err
    }

    fmt.Println("Item deleted successfully")
    return nil
}

// DelteDetailed will delete the item with the given pk and sk. This method should be
// prefered over the simpler Delete method
func (d *DynamoTable[T]) DeleteDetailed(pk, sk string) error {
    key := map[string]types.AttributeValue{
        d.pk: &types.AttributeValueMemberS{Value: pk},
        d.sk: &types.AttributeValueMemberS{Value: sk},
    }
    input := &dynamodb.DeleteItemInput{
        TableName: aws.String(d.tableName),
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

func (d *DynamoTable[T]) DeleteWithColumnName(pk, sk string) error {
    input := &dynamodb.DeleteItemInput{
        TableName: aws.String(d.tableName),
        Key: map[string]types.AttributeValue{
            d.pk: &types.AttributeValueMemberS{Value: pk},
            d.sk: &types.AttributeValueMemberS{Value: sk},
        },
    }

    _, err := d.svc.DeleteItem(context.TODO(), input)
    return err
}

func (d *DynamoTable[T]) DeleteGSIItem(gsiName, gsiPartitionKey, gsiPartitionValue  string) error {

    log.Println("gsiName: ", gsiName)
    log.Println("gsiPartitionKey: ", gsiPartitionKey)
    log.Println("gsiPartitionValue: ", gsiPartitionValue)

    // Query GSI to get primary key
    queryOut, err := d.svc.Query(context.TODO(), &dynamodb.QueryInput{
        TableName:              aws.String(d.tableName),
        IndexName:              aws.String(gsiName), // Specify the GSI name
        KeyConditionExpression: aws.String(fmt.Sprintf("%s = :v", gsiPartitionKey)),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":v": &types.AttributeValueMemberS{Value: gsiPartitionValue},
        },
    })
    if err != nil {
        log.Fatalf("Failed to query GSI: %v", err)
        return err
    }

    // Check if we found any items
    if len(queryOut.Items) == 0 {
        log.Println("No matching item found in GSI.")
        return nil
    }

    // Extract the primary key (assumes simple partition key and sort key structure)
    primaryKey := d.pk

    for _, item := range queryOut.Items {
        // Extract primary key values
        primaryKeyValue := item[primaryKey].(*types.AttributeValueMemberS).Value

        // Delete the item using primary key
        deleteInput := &dynamodb.DeleteItemInput{
            TableName: aws.String(d.tableName),
            Key: map[string]types.AttributeValue{
                primaryKey: &types.AttributeValueMemberS{Value: primaryKeyValue},
            },
        }

        _, err := d.svc.DeleteItem(context.TODO(), deleteInput)
        if err != nil {
            log.Fatalf("Failed to delete item: %v", err)
        }

        log.Printf("Deleted item with primary key: %s", primaryKeyValue)
    }

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

package models

import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DatabaseReadWriter[T any] interface {
    Read(id string) ([]T, error)
    ReadSpecific(pk , sk string) (T, error)

    Update(pk, sk string) error

    Create(t T, converter func (t T) map[string]types.AttributeValue)error

    Delete(pk string) error
    DeleteDetailed(pk , sk string) error
}

package models

import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DatabaseReadWriter[T any] interface {
    Read(id string) ([]T, error)
    ReadAll() ([]T, error)
    ReadSpecific(pk , sk string) (T, error)
    ReadCol(col, val string) ([]T, error)
    BatchRead(pks, sks []string) ([]T, error)
    BatchReadPks(pks []string) ([]T, error)

    Update(pk, sk, col, val string) error
    UpdateAll(pk, col, val string) error

    Create(t T, converter func (t T) map[string]types.AttributeValue)error

    Delete(pk string) error
    DeleteDetailed(pk , sk string) error
}

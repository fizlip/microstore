/// package models contains definitions of common types and models used accross the API.
package models

import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DatabaseReadWriter is an interface that defines all functions 
// that a database read writer has to implement
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

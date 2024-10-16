package models

type Table[T any] interface {
    Read(id string) ([]T, error)
    ReadSpecific(pk , sk string) (T, error)

    Update(pk, sk string) error

    Create(t T) error

    Delete(pk string) error
    DeleteDetailed(pk , sk string) error
}

package services

import (
    "fmt"

    "github.com/fizlip/microstore/internal/models"
)

func (t *TokenTable) Insert(t models.Token) error {
    fmt.Println("Insert token") 
}

func (t *TokenTable) Delete(username string) error {
    fmt.Println("delete token") 
}

func (t *TokenTable) Read(token string) error {
    fmt.Println("read token")
}

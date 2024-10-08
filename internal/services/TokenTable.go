package services

import (
    log "github.com/sirupsen/logrus"
)

type Token struct {
    Username    string
    Token       string
    ValidUntil  int 
    Provider    string
}

func (table *Token) Insert() error {
    log.Println("Insert token") 
    return nil
}

func (table *Token) Delete() error {
    log.Println("delete token") 
    return nil
}

func (table *Token) Read() error {
    log.Println("read token")
    return nil
}

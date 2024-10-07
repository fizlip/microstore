package models

type Token struct {
    Username    string
    Token       string
    ValidUntil  int 
}

type TokenTable interface {
  Insert(t Token) error
  Delete(username string) error
  Read(token string) error
}

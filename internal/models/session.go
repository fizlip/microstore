package models

type Session struct {
    UID         string
    F4rmer      string
    ChatHistory []Message 
} 

type Message struct {
    ID          int
    Content     string
    Role        string
}

// package models contains definitions of common types and models used accross the API.
package models

// Session is used to keep memory between chat session, i.e. users can keep chatting
// even if they update site or start using another device
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

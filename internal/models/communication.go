package models

type Communication struct {
	User    string `json:"user" dynamodbav:"user"`
	ID      string `json:"id" dynamodbav:"id"`
	Type    string `json:"type" dynamodbav:"type"`
	Content string `json:"content" dynamodbav:"content"`
	Number  int    `json:"number" dynamodbav:"number"`
}

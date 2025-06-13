package models

type Token struct {
	Username string `dynamodbav:"username" json:"username"`
	Server   string `dynamodbav:"server" json:"server"`
	Token    string `dynamodbav:"token" json:"token"`
	Provider string `dynamodbav:"provider" json:"provider"`
}

package models

type Token struct {
	Username string `dynamodbav:"username" json:"username"`
	Server   string `dynamodbav:"server" json:"server"`
	Token    Secret `dynamodbav:"token" json:"token"`
	Provider string `dynamodbav:"provider" json:"provider"`
}

type Secret struct {
	EncryptedData string `dynamodbav:"encryptedData" json:"encryptedData"`
	IV            string `dynamodbav:"iv" json:"iv"`
	Salt          string `dynamodbav:"salt" json:"salt"`
}

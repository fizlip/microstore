package models 

type ToolName struct {
    UTI     string      `dynamodbav:"uti" json:"uti"`
    IPs     []string    `dynamodbav:"ips" json:"ips"`
}

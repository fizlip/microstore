package models

type Customer struct {
	Email           string
	PaymentExecutor PaymentExecutor
	Billed          bool
	Verified        bool
}

type CustomerAnalytics struct {
    CustomerID      string              `json:"customer_id" dynamodbav:"customer_id"`
    Timestamp       string              `json:"timestamp" dynamodbav:"timestamp"`
    PaidRequests    int                 `json:"paid_requests" dynamodbav:"paid_requests"`
    FreeRequests    int                 `json:"free_requests" dynamodbav:"free_requests"`
    Products        map[string]ToolData `json:"products" dynamodbav:"products"`
}

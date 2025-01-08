// package models contains definitions of common types and models used accross the 
// API.
package models

// Customer is the model used to define a user of the system. Usually a customer
// will consume the products of a f4rmer
type Customer struct {
	Email           string
	PaymentExecutor PaymentExecutor
	Billed          bool
	Verified        bool
}

// CustomerAnalytics tracks customer usage of f4rmers which helps in calculating
// invoices.
type CustomerAnalytics struct {
    CustomerID      string              `json:"customer_id" dynamodbav:"customer_id"`
    Timestamp       string              `json:"timestamp" dynamodbav:"timestamp"`
    PaidRequests    int                 `json:"paid_requests" dynamodbav:"paid_requests"`
    FreeRequests    int                 `json:"free_requests" dynamodbav:"free_requests"`
    Products        map[string]ToolData `json:"products" dynamodbav:"products"`
}

package models 

type UserAnalytics struct {
    CustomerID      string                  `json:"customer_id" dynamodbav:"customer_id"`
    Timestamp       string                  `json:"timestamp" dynamodbav:"timestamp"`
    PaidRequests    int                     `json:"paid_requests" dynamodbav:"paid_requests"`
    FreeRequests    int                     `json:"free_requests" dynamodbav:"free_requests"`
    Products        map[string]ProductData  `json:"products" dynamodbav:"products"`
}

type ProductData struct {
    ProductID   string  `json:"product_id"`
    ProductName string  `json:"product_name"`
    Price       float64 `json:"price"`
    NumRequests int     `json:"num_requests"`
}

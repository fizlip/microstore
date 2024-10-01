package api

type UserAnalytics struct {
  CustomerID	string			`json:"customer_id"`
  Timestamp 	string			`json:"timestamp"`	
  PaidRequests	int64			`json:"paid_requests"`
  FreeRequests	int64			`json:"free_requests"`
  Products	map[string]ProductData	`json:"products"`
}

type ProductAnalytics struct {
  ProductID	string	`json:"product_id"`
  Timestamp	string	`json:"timestamp"`
  Requests 	uint64 	`json:"requests"`
  ProductPrice	float64 `json:"product_price"`
  Status	string	`json:"status"`
}

type ProductData struct {
  ProductID	string	`json:"product_id"`	
  ProductName	string	`json:"product_name"`
  Price		float64	`json:"price"`
  NumRequests	int64	`json:"num_requests"`
}

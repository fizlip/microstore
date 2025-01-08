// package models contains definitions of common types and models used accross the API.
package models

// ProductAnalytics is used to track the usage of tools
type ProductAnalytics struct {
	ProductID string    `json:"product_id" dynamodbav:"product_id"`
	Timestamp string    `json:"timestamp" dynamodbav:"timestamp"`
	Status    string    `json:"status" dynamodbav:"status"`
	Consumer  string    `json:"consumer" dynamodbav:"consumer"`
}

type IProductAnalyticsTable interface {
	Insert(p *Tool) error
	Read(pid string) ([]ProductAnalytics, error)
	Update(id string) (bool, error)
	Delete(id string) (bool, error)
}

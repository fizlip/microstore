package models

type ProductAnalytics struct {
	ProductID string `json:"product_id"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Consumer  string `json:"consumer"`
}

type IProductAnalyticsTable interface {
	Insert(p *Tool) error
	Read(pid string) ([]ProductAnalytics, error)
	Update(id string) (bool, error)
	Delete(id string) (bool, error)
}

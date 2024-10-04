package api

import (
	log "github.com/sirupsen/logrus"
)

type IProductAnalyticsTable interface {
	Insert(p *Product) error
	Read(pid string) ([]ProductAnalytics, error)
	Update(id string) (bool, error)
	Delete(id string) (bool, error)
}

type DynamoDBProductAnalyticsTable struct {
	TableName string
}

func (t *DynamoDBProductAnalyticsTable) Insert(p *Product) error {
	log.Println("Inserting analytics into table")

	pa := ProductAnalytics{
		ProductID: p.UID,
		Timestamp: "2024-10-01",
		Status:    "200",
		Consumer:  "test1233@mail.com",
	}

	mockProductAnalytics = append(mockProductAnalytics, pa)
	return nil
}

func (r *DynamoDBProductAnalyticsTable) Read(pid string) ([]ProductAnalytics, error) {
	log.Println("Reading reviews")
	res := []ProductAnalytics{}
	for _, value := range mockProductAnalytics {
		if value.ProductID == pid {
			res = append(res, value)
		}
	}
	return res, nil
}

func (r *DynamoDBProductAnalyticsTable) Update(pid string) (bool, error) {
	log.Println("Updating review")
	return true, nil
}

func (r *DynamoDBProductAnalyticsTable) Delete(pid string) (bool, error) {
	log.Println("Deleting review")
	return true, nil
}

func NewProductAnalyticsTable(tableName string) (*IProductAnalyticsTable, error) {
	var analyticsTable IProductAnalyticsTable = &DynamoDBProductAnalyticsTable{tableName}
	return &analyticsTable, nil
}

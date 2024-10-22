package services 

import (
	log "github.com/sirupsen/logrus"
    "github.com/fizlip/microstore/internal/models"
)

type DynamoDBProductAnalyticsTable struct {
	TableName string
}

type T DynamoDBProductAnalyticsTable

func (t *T) Insert(p *models.Tool) error {
	log.Println("Inserting analytics into table")

	//pa := models.ProductAnalytics{
	//	ProductID: p.UID,
	//	Timestamp: "2024-10-01",
	//	Status:    "200",
	//	Consumer:  "test1233@mail.com",
	//}

    // Todo: implement interface for pa table
	// mocks.MockProductAnalytics = append(mocks.MockProductAnalytics, pa)
	return nil
}

func (r *T) Read(pid string) ([]models.ProductAnalytics, error) {
	log.Println("Reading reviews")
	//res := []models.ProductAnalytics{}
	//for _, value := range mocks.MockProductAnalytics {
	//	if value.ProductID == pid {
	//		res = append(res, value)
	//	}
	//}
	return []models.ProductAnalytics{}, nil
}

func (r *T) Update(pid string) (bool, error) {
	log.Println("Updating review")
	return true, nil
}

func (r *T) Delete(pid string) (bool, error) {
	log.Println("Deleting review")
	return true, nil
}

func NewProductAnalyticsTable(tableName string) (*models.IProductAnalyticsTable, error) {
	var analyticsTable models.IProductAnalyticsTable = &T{tableName}
	return &analyticsTable, nil
}

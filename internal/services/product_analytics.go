// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
	"github.com/fizlip/microstore/internal/models"
	log "github.com/sirupsen/logrus"
)

// ProductAnalyticsService is only used to do CRUD operations on the
// f4-analytics table.
type ProductAnalyticsService struct {
	analyticStore models.DatabaseReadWriter[models.ProductAnalytics]
}

func (p *ProductAnalyticsService) GetRecentCalls(uti string) ([]models.ProductAnalytics, error) {
	return p.analyticStore.Read(uti)
}

func (p *ProductAnalyticsService) Insert(t *models.Tool) error {
	log.Println("Inserting analytics into table")
	return nil
}

func (p *ProductAnalyticsService) Read(pid string) ([]models.ProductAnalytics, error) {
	log.Println("Reading reviews")
	//res := []models.ProductAnalytics{}
	//for _, value := range mocks.MockProductAnalytics {
	//	if value.ProductID == pid {
	//		res = append(res, value)
	//	}
	//}
	return []models.ProductAnalytics{}, nil
}

func (p *ProductAnalyticsService) Update(pid string) (bool, error) {
	log.Println("Updating review")
	return true, nil
}

func (p *ProductAnalyticsService) Delete(pid string) (bool, error) {
	log.Println("Deleting review")
	return true, nil
}

func NewActionAnalyticsService() (*ProductAnalyticsService, error) {
	analyticsT, err := models.NewDynamoTable[models.ProductAnalytics]("f4-product-analytics", "product_id", "timestamp")
	if err != nil {
		return &ProductAnalyticsService{}, err
	}
	productAnalytics := ProductAnalyticsService{analyticStore: analyticsT}
	return &productAnalytics, nil
}

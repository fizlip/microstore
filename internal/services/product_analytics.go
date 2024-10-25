package services 

import (
	log "github.com/sirupsen/logrus"
    "github.com/fizlip/microstore/internal/models"
)

type ProductAnalyticsService struct {
    AnalyticStorage     models.DatabaseReadWriter[models.ProductAnalytics]
}

func (p *ProductAnalyticsService) Insert(t *models.Tool) error {
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

func NewProductAnalyticsService(tableName string) (*ProductAnalyticsService, error) {
    analyticT, err := models.NewMockDynamoTable[models.ProductAnalytics]("f4-analytics", "developer", "uid")
    if err != nil {
        return &ProductAnalyticsService{}, err 
    }
    productAnalytics := ProductAnalyticsService{AnalyticStorage: analyticT}
	return &productAnalytics, nil
}

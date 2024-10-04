package services 

import (
	log "github.com/sirupsen/logrus"
    "github.com/fizlip/microstore/api/models"
)

type Service struct {
    Product models.Product
}

func (s *Service) Execute(cid string) (interface{}, error) {
	log.Println("Log new execution entry for: " + s.Product.UID)
	s.CreateLog()
	return "Hello world", nil
}

func (s *Service) CreateLog() error {
	table, _ := NewProductAnalyticsTable("product_analytics")
	err := (*table).Insert(&s.Product)
	if err != nil {
		log.Println("Could not create new log")
	}
	return nil
}

func NewProduct(uid string, price float64) (*Service, error) {
	return &Service{Product: models.Product{UID: uid, Price: price}}, nil
}

package services

import (
    "github.com/fizlip/microstore/internal/models"

    // log "github.com/sirupsen/logrus"
)

type Customer struct {
    Name            string
	Email           string
    analyticsStore  models.DatabaseReadWriter[models.ProductAnalytics]
}

func (c *Customer) GetAnalytics() ([]models.ProductAnalytics, error) {
    return []models.ProductAnalytics{}, nil
}

func (c *Customer) GetUsage() ([]models.ProductAnalytics, error) {
    return []models.ProductAnalytics{}, nil
}

func NewCustomer(username string) (Customer, error) {
    return Customer{}, nil
}

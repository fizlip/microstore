package mocks 

import (
    "github.com/fizlip/microstore/api/models"
)

var MockProductAnalytics = []models.ProductAnalytics{
	{
		ProductID: "p_001",
		Timestamp: "2024-10-01",
		Status:    "200",
		Consumer:  "test@mail.com",
	},
}

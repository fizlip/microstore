package mocks 

type ProductAnalytics struct {
    ProductID   string
    Timestamp   string
    Status      string
    Consumer    string
}

var MockProductAnalytics = []ProductAnalytics{
	{
		ProductID: "p_001",
		Timestamp: "2024-10-01",
		Status:    "200",
		Consumer:  "test@mail.com",
	},
}

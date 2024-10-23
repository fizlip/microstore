package models

import (
)

type Customer struct {
	Email           string
	PaymentExecutor PaymentExecutor
	Billed          bool
	Verified        bool
}

type CustomerAnalytics struct {
    CustomerID      string
    Timestamp       string
    PaidRequests    int
    FreeRequests    int
    Products        map[string]ToolData
}

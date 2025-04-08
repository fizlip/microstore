package models

type Invoice struct {
    InvoiceID           string
    BillingPeriod       string
    Name                string
    RequestsMade        int64
    ActionsUsed         []string
    ActionPrices        []float64
    TotalCostExTax     float64
    TotalCostIncTax    float64
}

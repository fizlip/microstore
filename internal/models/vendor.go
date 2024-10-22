package models

type Vendor struct {
    Name            string
	Email           string
	Rating          float64
	ProductsCreated int
	RequestHandled  int64
	ProductCatalog  []Tool
	PaymentExecutor PaymentExecutor
    ToolStorage     DatabaseReadWriter[Tool]
}

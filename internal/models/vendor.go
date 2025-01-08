// package models contains definitions of common types and models used accross the API.
package models

// Vendor is a user that also creates tools.
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

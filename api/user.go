package api 

import (
  "fmt"
)

type User interface {
  Register()
  Login()
  Logout()
}

type Customer struct {
  Email 		string
  PaymentExecutor 	PaymentExecutor
  Billed		bool
  Verified 		bool
}

func (c *Customer) Login() {
  fmt.Println("Login: "  + c.Email)
}

func (c *Customer) Register() {
  fmt.Println("Register: "  + c.Email)
}

func (c *Customer) Logout() {
  fmt.Println("Logout: "  + c.Email)
}

func (c *Customer) RunProduct(ev *F4SubscriptionEvent) (*PaymentResponse, error) {
  return c.PaymentExecutor.CreateMeterEvent(ev)
}

func (c *Customer) WriteReview() {
  fmt.Println("Write review: "  + c.Email)
}

func (c *Customer) RateReview() {
  fmt.Println("rate review: "  + c.Email)
}

func (c *Customer) GetUsage() {
  fmt.Println("get usage: "  + c.Email)
}

func NewCustomer(email string) (*Customer, error){
  var mockPex PaymentExecutor = &mockStripePaymentExecutor{PK: "123"}

  var customer *Customer = &Customer{Email: email, PaymentExecutor: mockPex}

  return customer, nil
}

type Vendor struct {
  Email			string
  Rating		float64
  ProductsCreated	int
  RequestHandled	int64
  ProductCatalog	[]Product
  PaymentExecutor	PaymentExecutor
}

func (c *Vendor) Login() {
  fmt.Println("Login: "  + c.Email)
}

func (c *Vendor) Register() {
  fmt.Println("Register: "  + c.Email)
}

func (c *Vendor) Logout() {
  fmt.Println("Logout: "  + c.Email)
}

func (c *Vendor) CreateProduct(p *Product) {
  fmt.Println("Product created by: "  + c.Email)
}

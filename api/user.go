package api 

import (
  "fmt"
  "github.com/stripe/stripe-go/v79"
  log "github.com/sirupsen/logrus"
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

  product, _ := NewProduct("some_random_id")
  product.Execute(c.Email)

  return c.PaymentExecutor.CreateMeterEvent(ev)
}

func (c *Customer) WriteReview() {
  fmt.Println("Write review: "  + c.Email)
}

func (c *Customer) RateReview() {
  fmt.Println("rate review: "  + c.Email)
}

func (c *Customer) GetUsage() (*stripe.Invoice, error){
  fmt.Println("get usage: "  + c.Email)
  return c.PaymentExecutor.GetUpcomingInvoice("cus_QsQCZFMdRmYIP8")
}

func NewCustomer(email string) (*Customer, error){
  var mockPex PaymentExecutor = &mockStripePaymentExecutor{
    PK: "123", 
    SK: "sk_test_51Q0fHUG9IK8VP5Ow1GeRlhjEG6fQPXEozXyegjNbNTl7QT1uCZGcz6rCKcXBGMWkWNJGDQyR88bfIGjbP53NzGBs00uYLlMnJp"}

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

func (c *Vendor) ExecuteProduct(cid string) {
  log.Println("Give moeny to vendor")
} 

func VendorFromProductID(productID string) (*Vendor, error) {
  var mockPex PaymentExecutor = &mockStripePaymentExecutor{
    PK: "123", 
    SK: "sk_test_51Q0fHUG9IK8VP5Ow1GeRlhjEG6fQPXEozXyegjNbNTl7QT1uCZGcz6rCKcXBGMWkWNJGDQyR88bfIGjbP53NzGBs00uYLlMnJp"}


  return &Vendor{
    Email: "vendor@example.com",
    Rating: 100.0,
    ProductsCreated: 120,
    RequestHandled: 10, 
    ProductCatalog: []Product{},
    PaymentExecutor: mockPex}, nil
}

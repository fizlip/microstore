package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v79"
    "github.com/fizlip/microstore/api/mocks"
    "github.com/fizlip/microstore/api/types"
)

type User interface {
	Register()
	Login()
	Logout()
}

type Customer struct {
	Email           string
	PaymentExecutor PaymentExecutor
	Billed          bool
	Verified        bool
}

func (c *Customer) Login() {
	fmt.Println("Login: " + c.Email)
}

func (c *Customer) Register() {
	fmt.Println("Register: " + c.Email)
}

func (c *Customer) Logout() {
	fmt.Println("Logout: " + c.Email)
}

// RunProduct will execute the product in the request and create
// UserAnalytics and ProductAnalytics
func (c *Customer) RunProduct(ev *F4SubscriptionEvent) (*PaymentResponse, error) {
	// Create product object
	product, _ := NewProduct(ev.ProductID, 0.01)
	// Execute product
	product.Execute(c.Email)
	// If successfull create new analytics (this is used to calculate user invoice)
	c.CreateUserAnalytic(false, product)
	// Then create product analytics (this is used as product logs and to calculate vendor pay)

	return c.PaymentExecutor.CreateMeterEvent(ev)
}

func (c *Customer) WriteReview() {
	fmt.Println("Write review: " + c.Email)
}

func (c *Customer) RateReview() {
	fmt.Println("rate review: " + c.Email)
}

func (c *Customer) GetUsage() (*stripe.Invoice, error) {
	fmt.Println("get usage: " + c.Email)
	return c.PaymentExecutor.GetUpcomingInvoice("cus_QsQCZFMdRmYIP8")
}

func (c *Customer) GetAnalytics() ([]types.UserAnalytics, error) {
	log.Println("Getting analytics for: " + c.Email)
	return mocks.MockUserAnalytics, nil
}

func (c *Customer) CreateUserAnalytic(paid bool, product *Product) {
	var _paid int64 = 0
	var free int64 = 0

	if product.Price > 0.0 {
		_paid = 1
	} else {
		free = 1
	}

	today := "2024-09-25"
	userData := mocks.MockUserAnalytics

	for i, elem := range userData {
		if elem.Timestamp == today && elem.CustomerID == c.Email {
			if free > 0 {
				userData[i].FreeRequests += 1
			} else {
				userData[i].PaidRequests += 1
			}

			entry, exists := userData[i].Products[product.UID]
			if exists {
				entry.NumRequests += 1
				userData[i].Products[product.UID] = entry
			} else {
				userData[i].Products[product.UID] = types.ProductData{ProductID: product.UID, ProductName: product.Title, Price: product.Price, NumRequests: 1}
			}

			mocks.MockUserAnalytics = userData
			return
		}
	}

	newAnalytic := types.UserAnalytics{
		CustomerID:   c.Email,
		Timestamp:    today,
		PaidRequests: _paid,
		FreeRequests: free,
		Products: map[string]types.ProductData{
			product.UID: {ProductID: product.UID, ProductName: product.Title, Price: product.Price, NumRequests: 1},
		},
	}
	mocks.MockUserAnalytics = append(userData, newAnalytic)
	log.Println("Creating new analytic")
}

func NewCustomer(email string) (*Customer, error) {
	var mockPex PaymentExecutor = &mockStripePaymentExecutor{
		PK: "123",
		SK: "sk_test_51Q0fHUG9IK8VP5Ow1GeRlhjEG6fQPXEozXyegjNbNTl7QT1uCZGcz6rCKcXBGMWkWNJGDQyR88bfIGjbP53NzGBs00uYLlMnJp"}

	var customer *Customer = &Customer{Email: email, PaymentExecutor: mockPex}

	return customer, nil
}

type Vendor struct {
	Email           string
	Rating          float64
	ProductsCreated int
	RequestHandled  int64
	ProductCatalog  []Product
	PaymentExecutor PaymentExecutor
}

func (c *Vendor) Login() {
	fmt.Println("Login: " + c.Email)
}

func (c *Vendor) Register() {
	fmt.Println("Register: " + c.Email)
}

func (c *Vendor) Logout() {
	fmt.Println("Logout: " + c.Email)
}

func (c *Vendor) CreateProduct(p *Product) {
	fmt.Println("Product created by: " + c.Email)
}

func (c *Vendor) ExecuteProduct(cid string) {
	log.Println("Give moeny to vendor")
}

func (c *Vendor) GetProducts() []Product {
	return c.ProductCatalog
}

func (c *Vendor) GetProductUsage() ([]ProductAnalytics, error) {
	vendorProducts := c.GetProducts()
	pa, err := NewProductAnalyticsTable("product_analytics")
	if err != nil {
		return []ProductAnalytics{}, err
	}

	res := []ProductAnalytics{}

	for _, product := range vendorProducts {
		fmt.Println(product.UID)
		analytics, err := (*pa).Read(product.UID)
		if err != nil {
			log.Error("Cuould not read", product.UID)
		} else {
			res = append(res, analytics...)
		}
	}
	return res, nil
}

func NewVendor(uid string) (*Vendor, error) {
	var mockPex PaymentExecutor = &mockStripePaymentExecutor{
		PK: "123",
		SK: "sk_test_51Q0fHUG9IK8VP5Ow1GeRlhjEG6fQPXEozXyegjNbNTl7QT1uCZGcz6rCKcXBGMWkWNJGDQyR88bfIGjbP53NzGBs00uYLlMnJp"}

	return &Vendor{
		Email:           "vendor@example.com",
		Rating:          100.0,
		ProductsCreated: 2,
		RequestHandled:  10,
		ProductCatalog:  mockProductCatalog,
		PaymentExecutor: mockPex}, nil
}

func VendorFromProductID(productID string) (*Vendor, error) {
	var mockPex PaymentExecutor = &mockStripePaymentExecutor{
		PK: "123",
		SK: "sk_test_51Q0fHUG9IK8VP5Ow1GeRlhjEG6fQPXEozXyegjNbNTl7QT1uCZGcz6rCKcXBGMWkWNJGDQyR88bfIGjbP53NzGBs00uYLlMnJp"}

	return &Vendor{
		Email:           "vendor@example.com",
		Rating:          100.0,
		ProductsCreated: 120,
		RequestHandled:  10,
		ProductCatalog:  []Product{},
		PaymentExecutor: mockPex}, nil
}

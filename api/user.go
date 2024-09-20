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
  UID 	string
  Email string
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

func (c *Customer) RunProduct() {
  fmt.Println("Logout: "  + c.Email)
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

type Vendor struct {
  Email			string
  Rating		float64
  ProductsCreated	int
  RequestHandled	int64
  ProductCatalog	[]Product
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

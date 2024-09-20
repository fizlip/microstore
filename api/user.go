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

type Vendor struct {
  Rating		float64
  ProductsCreated	int
  RequestHandled	int64
  ProductCatalog	[]Product
}

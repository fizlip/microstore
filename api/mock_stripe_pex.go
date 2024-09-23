package api 

import (
  "fmt"
  log "github.com/sirupsen/logrus"
)

type mockStripePaymentExecutor struct {
  PK	string
}

func (p *mockStripePaymentExecutor) CreateMeterEvent() {
  log.Println("Create a stripe test meter event")
}

func (p *mockStripePaymentExecutor) GetUpcomingInvoice() {
  fmt.Println("Get an upcoming invoice")
}

func (p *mockStripePaymentExecutor) GetUpcomingPayment(uid string) {
  fmt.Println("Get an upcoming payment")
}

func (p *mockStripePaymentExecutor) GetProductUsage() {
  fmt.Println("Product usage")
}

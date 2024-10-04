package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type StripePaymentExecutor struct {
	PK string
}

func (p *StripePaymentExecutor) CreateMeterEvent() (*PaymentResponse, error) {
	log.Println("Create a stripe meter event")
	return &PaymentResponse{}, nil
}

func (p *StripePaymentExecutor) GetUpcomingInvoice() {
	fmt.Println("Get an upcoming invoice")
}

func (p *StripePaymentExecutor) GetUpcomingPayment(uid string) {
	fmt.Println("Get an upcoming payment")
}

func (p *StripePaymentExecutor) GetProductUsage() {
	fmt.Println("Product usage")
}

package api 

import (
  "github.com/google/uuid"
  "github.com/stripe/stripe-go/v79"
)

type F4SubscriptionEvent struct {
  EventID	uuid.UUID	
  CustomerID	string
  Value		string
  ProductID	string
}

type PaymentResponse struct {
  UID		uuid.UUID	
  Status	string
  ProductID	string
  Value		string
  Message	string
}

type PaymentExecutor interface {
  CreateMeterEvent(ev *F4SubscriptionEvent) (*PaymentResponse, error)
  GetUpcomingInvoice(cid string) (*stripe.Invoice, error)
  GetUpcomingPayment(uid string)
  GetProductUsage()
}

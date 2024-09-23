package api 

import (
  "fmt"
  log "github.com/sirupsen/logrus"
  stripe "github.com/stripe/stripe-go/v79"
  "github.com/stripe/stripe-go/v79/billing/meterevent"
)

type mockStripePaymentExecutor struct {
  PK	string
}

func (p *mockStripePaymentExecutor) CreateMeterEvent(ev *F4SubscriptionEvent) (*PaymentResponse, error) {

  log.Println("Create a stripe test meter event")

  stripe.Key = "sk_test_51Q0fHUG9IK8VP5Ow1GeRlhjEG6fQPXEozXyegjNbNTl7QT1uCZGcz6rCKcXBGMWkWNJGDQyR88bfIGjbP53NzGBs00uYLlMnJp"
  params := &stripe.BillingMeterEventParams{
    EventName: stripe.String(ev.ProductID),
    Payload: map[string]string{"value": ev.Value, "stripe_customer_id": ev.CustomerID},
  };
  result, err := meterevent.New(params);
  if err != nil {
    log.Error(err)
    return &PaymentResponse{}, err
  }

  log.Println("New event for: " + result.EventName)
  return &PaymentResponse{UID: "123", Status: "200", ProductID: ev.ProductID, Value: ev.Value, Message: "Payment was successful"}, nil
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

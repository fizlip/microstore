package api

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	stripe "github.com/stripe/stripe-go/v79"

	"github.com/stripe/stripe-go/v79/invoice"
)

type mockStripePaymentExecutor struct {
	PK string
	SK string
}

// CreateMeterEvent registers a new payemtn for a products, it is run using
// stripe subscription events
//func (p *mockStripePaymentExecutor) CreateMeterEvent(ev *F4SubscriptionEvent) (*PaymentResponse, error) {
//
//	log.Println("Creating a stripe test meter event")
//
//	// Secret key stored in pex
//	stripe.Key = p.SK
//	params := &stripe.BillingMeterEventParams{
//		EventName: stripe.String(ev.ProductID),
//		Payload:   map[string]string{"value": ev.Value, "stripe_customer_id": ev.CustomerID},
//	}
//	result, err := meterevent.New(params)
//	if err != nil {
//		log.Error(err)
//		// Generic error message if payment is not sucessful
//		return &PaymentResponse{UID: uuid.New(), Status: "400", Message: "Could process create payment"}, err
//	}
//
//	log.Println("New event for: " + result.EventName)
//	// Success message if no error occured
//	return &PaymentResponse{
//		UID:       uuid.New(),
//		Status:    "200",
//		ProductID: ev.ProductID,
//		Value:     ev.Value,
//		Message:   "Payment was successful"}, nil
//}

func (p *mockStripePaymentExecutor) GetUpcomingInvoice(cid string) (*stripe.Invoice, error) {
	stripe.Key = p.SK

	fmt.Println("Get an upcoming invoice")

	params := &stripe.InvoiceUpcomingParams{Customer: stripe.String(cid)}
	result, err := invoice.Upcoming(params)
	if err != nil {
		log.Error("Could not get invoice")
		return result, err
	}
	log.Println("Generated invoice for: " + cid)
	return result, nil

}

func (p *mockStripePaymentExecutor) GetUpcomingPayment(uid string) {
	fmt.Println("Get an upcoming payment")
}

func (p *mockStripePaymentExecutor) GetProductUsage() {
	fmt.Println("Product usage")
}

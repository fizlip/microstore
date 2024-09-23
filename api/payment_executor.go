package api 

type F4SubscriptionEvent struct {
  CustomerID	string
  Value		string
  ProductID	string
}

type PaymentResponse struct {
  UID		string
  Status	string
  ProductID	string
  Value		string
  Message	string
}

type PaymentExecutor interface {
  CreateMeterEvent(ev *F4SubscriptionEvent) (*PaymentResponse, error)
  GetUpcomingInvoice()
  GetUpcomingPayment(uid string)
  GetProductUsage()
}

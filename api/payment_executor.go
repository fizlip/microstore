package api 

type PaymentExecutor interface {
  CreateMeterEvent()
  GetUpcomingInvoice()
  GetUpcomingPayment(uid string)
  GetProductUsage()
}

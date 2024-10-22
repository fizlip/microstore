package models

import (
)

type Customer struct {
	Email           string
	PaymentExecutor PaymentExecutor
	Billed          bool
	Verified        bool
}

package userhandlers

import (
	"encoding/json"
	"net/http"

    "github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

// GetInvoice returns an invoice object from the payment provider that is used
func GetInvoice(w http.ResponseWriter, r *http.Request) {
	var err error

    customer_id := r.Header.Get("X-Username")
    user, err := services.NewCustomer(customer_id)
	if err != nil {
        log.Error("Cannot connect to customer service", err)
        api.InternalErrorHandler(w)
        return
	}

	// And run the product with it
	result, err := (*user).GetInvoice()
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

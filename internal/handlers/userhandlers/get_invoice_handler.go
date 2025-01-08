package userhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

// GetInvoice returns an invoice object from the payment provider that is used
func GetInvoice(w http.ResponseWriter, r *http.Request) {
	var err error
	//user, err := api.NewCustomer("some_email")
	//if err != nil {
    //    log.Error(err)
	//	api.InternalErrorHandler(w)
	//	return
	//}

	// And run the product with it
	//result, err := (*user).GetUsage()
	err = json.NewEncoder(w).Encode("")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

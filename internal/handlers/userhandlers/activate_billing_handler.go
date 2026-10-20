package userhandlers

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/stripe/stripe-go/v79"
	"net/http"

	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

func ActivateBilling(w http.ResponseWriter, r *http.Request) {
	var params = api.UsernameParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	pk := ""

	//_, err := api.NewCustomer("some_email")
	//if err != nil {
	//	panic(err)
	//}

	err := decoder.Decode(&params, r.URL.Query())

	stripe.Key = pk
	//c_params := &stripe.CustomerParams{
	//	Name:  stripe.String("Jenny Rosen"),
	//	Email: stripe.String("jennyrosen@example.com"),
	//}
	//result, err := customer.New(c_params)

	//if err != nil {
	//	panic(err)
	//}

	log.Println("Activate billing for: ", params.Username)

	// And run the product with it
	err = json.NewEncoder(w).Encode(params.Username)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

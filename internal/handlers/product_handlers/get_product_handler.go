// package product_handlers provides functions used to interact with 
// products.
package product_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/tools"
	"github.com/gorilla/schema"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetProduct will return one image with the UID given in the URL param.
// e.g. /product/get/image?uid=1234 will return image with the UID that is
// 1234. If it does not exist an empty string is returned.
// TODO: CHECK IF THIS IS NEEDED
func GetProduct(w http.ResponseWriter, r *http.Request) {
	var params = api.GetProductParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	var response = (*database).GetProduct(params.UID)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

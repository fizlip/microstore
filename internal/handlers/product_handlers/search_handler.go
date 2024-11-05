// package product_handlers provides functions used to interact with 
// products.
package product_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/tools"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Search will return the response to a search query as an array of 
// products.
func Search(w http.ResponseWriter, r *http.Request) {
	var params = api.SearchParams{}
	var err error

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		panic(err)
	}

	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	var response = (*database).SearchProduct(params.Query)

	log.Println("search product")

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

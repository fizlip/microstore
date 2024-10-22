package product_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/tools"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetAllImages will return one page of all images available in the f4mhouse OCI repo
// A correct response will return an array of Image values.
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

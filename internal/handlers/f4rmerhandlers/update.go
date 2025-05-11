// package f4rmerhandlers provides functions used to interact with f4rmers.
package f4rmerhandlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/fizlip/microstore/api"

	"github.com/fizlip/microstore/internal/models"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// Update will change the value of a field in an existing f4rmer, e.g. change the title.
// To update a f4rmer you need to provide f4rmer uid, field and the value that the field
// should be
//
// example: PATCH https://<hostname>/user/f4rmer/update
// response:
//
//	SUCCESS: 200, Update successfull
//	FAILURE: 403, An unexpected error occured
func Update(w http.ResponseWriter, r *http.Request) {
	var err error
	var response api.Success

	var instigator = r.Header.Get("X-Username")

	user, err := services.NewCustomer(instigator)

	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	var f4 models.F4rmer
	err = json.Unmarshal(body, &f4)
	if err != nil {
		// Could not parse object from body, return error
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	} else {
		// Add necessary fields to f4
		if f4.Creator != instigator {
			log.Error(err)
			api.AuthorizationErrorHandler(w, errors.New("Update not authorized. Only owner of f4rmers can update f4rmers."))
			return
		}

		err = (*user).CreateF4rmer(f4)
		response = api.Success{
			Code:    200,
			Message: "F4rmer was updated ",
		}
	}

	log.Println("F4rmer updated")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

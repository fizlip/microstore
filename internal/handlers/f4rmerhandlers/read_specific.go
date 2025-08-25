// package f4rmerhandlers provides functions used to interact with f4rmers.
package f4rmerhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// ReadSpecific will return a single f4rmer by their UID if it exists.
// To read a f4rmer you need to provide a username and uid (uinque ID that each f4rmer has)
//
// example: GET https://<hostname>/user/f4rmer/read_specific?username=<username>&uid=<f4rmer_UUID>
// response:
//
//	SUCCESS: 200, <f4rmer_object>
//	FAILURE: 403, An unexpected error occured
func ReadSpecific(w http.ResponseWriter, r *http.Request) {
	var err error
	var user *services.Customer
	usr := r.URL.Query().Get("username")
	uid := r.URL.Query().Get("uid")

	user, err = services.NewCustomer(usr)
	result, err := user.ReadF4rmer(usr, uid)

	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

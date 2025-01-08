// package f4rmer_handlers provides functions used to interact with f4rmers.
package f4rmer_handlers 

import (
	"encoding/json"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	"net/http"

	log "github.com/sirupsen/logrus"
)

// Read will return all f4rmers owned by a user.
// To read a f4rmer you need to be logged in with a valid token and provide the X-Username
// header 
//
// example: GET https://<hostname>/user/f4rmer/read
// response: 
//      SUCCESS: 200, <array_of_f4rmers>
//      FAILURE: 403, An unexpected error occured
//
func Read(w http.ResponseWriter, r *http.Request) {
    var err error
    var user *services.Customer
    var username = r.Header.Get("X-Username")

    user , err = services.NewCustomer(username)
    result, err := user.ReadF4rmers()
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

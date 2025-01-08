// package f4rmerhandlers provides functions used to interact with f4rmers.
package f4rmerhandlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Update will change the value of a field in an existing f4rmer, e.g. change the title.
// To update a f4rmer you need to provide f4rmer uid, field and the value that the field
// should be 
//
// example: UPDATE https://<hostname>/user/f4rmer/update
// response: 
//      SUCCESS: 200, Update successfull
//      FAILURE: 403, An unexpected error occured
func Update(w http.ResponseWriter, r *http.Request) {
    var err error
    err = json.NewEncoder(w).Encode("")
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

// package f4rmer_handlers provides functions used to interact with f4rmers.
package f4rmerhandlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Delete will delete a f4rmer from the database. To delete a f4rmer you need to
// provide the creator and the uid of the f4rmer. 
//
// example: DELETE https://<hostname>/user/f4rmer/delete?uid=<f4rmer_uuid>
// response: 
//      SUCCESS: 200, F4rmer has been deleted
//      FAILURE: 403, An unexpected error occured
func Delete(w http.ResponseWriter, r *http.Request) {
    var err error
    err = json.NewEncoder(w).Encode("")
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

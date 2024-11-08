// Package f4rmer_handlers.
package f4rmer_handlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	//"github.com/fizlip/microstore/internal/services"
	//"github.com/fizlip/microstore/internal/models"
	//"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ... 
func Update(w http.ResponseWriter, r *http.Request) {
    var err error
    err = json.NewEncoder(w).Encode("")
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

// Package f4rmer_handlers.
package f4rmer_handlers 

import (
	"encoding/json"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	"net/http"

	log "github.com/sirupsen/logrus"
)

// ... 
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

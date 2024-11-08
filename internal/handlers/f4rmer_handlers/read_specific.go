// Package f4rmer_handlers.
package f4rmer_handlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	//"github.com/fizlip/microstore/internal/models"
	//"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ... 
func ReadSpecific(w http.ResponseWriter, r *http.Request) {
    var err error
    var user *services.Customer
    usr := r.URL.Query().Get("username")
    uid := r.URL.Query().Get("uid")

    log.Println(usr)
    log.Println(uid)

    user , err = services.NewCustomer(usr)
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

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
func AddTool(w http.ResponseWriter, r *http.Request) {
    var err error
    var user *services.Customer
    usr := r.URL.Query().Get("username")
    uid := r.URL.Query().Get("uid")
    tid := r.URL.Query().Get("tid")

    user , err = services.NewCustomer(usr)
    err = user.AddTool(usr, uid, tid)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    err = json.NewEncoder(w).Encode("Tool was added")
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

// package f4rmer_handlers provides functions used to interact with f4rmers.
package f4rmerhandlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"
    "io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Run ...
func Run(w http.ResponseWriter, r *http.Request) {
    var err error
    var user *services.Customer
    var instigator = r.Header.Get("X-Username")

    // Read body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return 
    }
    var runParams api.RunBody
    err = json.Unmarshal(body, &runParams)

    user, err = services.NewCustomer(instigator)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
    var tns *services.ToolNameService
    tns, err = services.NewToolNameService()
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}
    result, err := (*tns).GetTool(runParams.UTI)
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

    res, err := user.Run(result[0].IPs[0], runParams.Endpoint, runParams.Params)
    if err != nil {
        log.Println("ERROR")
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    err = json.NewEncoder(w).Encode(res)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

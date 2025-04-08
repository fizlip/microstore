// package f4rmer_handlers provides functions used to interact with f4rmers.
package f4rmerhandlers 

import (
    "fmt"
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"
    "io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Run will execute the UTI and increment all necessary analytics metrics
func Run(w http.ResponseWriter, r *http.Request) {
    var err error

    var user *services.Customer
    var instigator = r.Header.Get("X-Username")

    // Read body, this will contain variables given to the function that's executed
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return 
    }
    // unmarshal body to get JSON
    var runParams api.RunBody
    err = json.Unmarshal(body, &runParams)

    user, err = services.NewCustomer(instigator)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
    // the tns is used to get the tools IP (like a DNS in classical web)
    var tns *services.ToolNameService
    tns, err = services.NewToolNameService()
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
    result, err := (*tns).GetTool(runParams.UTI)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

    // Execute the action 
    res, err := user.Run(result[0].IPs[0], runParams.Endpoint, runParams.Params, result[0].UTI)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }

    log.Println(fmt.Sprintf("%s was executed sucessfully", result[0].UTI))
    err = json.NewEncoder(w).Encode(api.RunSuccess{Code: 200, Message: res})
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

// package f4rmer_handlers provides functions used to interact with f4rmers.
package f4rmerhandlers 

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/internal/models"
    "github.com/google/uuid"
	"io/ioutil"
	"net/http"
    "time"

	log "github.com/sirupsen/logrus"
)

// Create will create a new f4rmer. To add a new f4rmer you need to provider
// a f4rmer object in the body of the request. the f4rmer struct is defined in
// models f4rmer.
//
// example: POST https://<hostname>/user/f4rmer/create 
//                -d {
//                      uid: string, 
//                      job_description:string
//                      toolbox:string
//                      creator:string
//                      created:string
//                   }
// response: 
//      SUCCESS: 200, F4rmer was successfully created
//      FAILURE: 403, An unexpected error occured
//
func Create(w http.ResponseWriter, r *http.Request) {
    var err error
    var response api.Success

    var instigator = r.Header.Get("X-Username")

    user, err := services.NewCustomer(instigator)
    
    // Read body
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return 
    }

    var f4 models.F4rmer
    err = json.Unmarshal(body, &f4)
    if err != nil {
        // Could not parse object from body, return error 
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    } else {
        // Add necessary fields to f4
        f4.Creator = instigator
        f4.UID = uuid.New().String() 
        f4.Created = time.Now().Format("2006-01-02 15:04:05") 
        err = (*user).CreateF4rmer(f4)

        response = api.Success{
            Code:    200,
            Message: "F4rmer was successfully created",
        }
    }

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

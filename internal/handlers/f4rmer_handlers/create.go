// Package f4rmer_handlers.
package f4rmer_handlers 

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

// ... 
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
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    } else {
        log.Println("F4rmer: ", f4)
        f4.Creator = instigator
        f4.UID = uuid.New().String() 
        f4.Created = time.Now().Format("2006-01-02 15:04:05") 
        err = (*user).CreateF4rmer(f4)

        response = api.Success{
            Code:    200,
            Message: "Review was successfully created",
        }
    }

    err = json.NewEncoder(w).Encode(response)
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

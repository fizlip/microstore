package user_handlers

import (
	"encoding/json"
	//"github.com/google/uuid"
	"net/http"

	"github.com/fizlip/microstore/internal/services"
	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

// RunProduct will execute a tool using the execution engine and update the
// usage of the tool that has been used
func RunProduct(w http.ResponseWriter, r *http.Request) {
	var err error
    var data map[string]interface{}

    user_id := r.Header.Get("X-Username")
    tool_id := r.URL.Query().Get("tid")
    // TODO: Get this from body
    datastr := `{
        "name": "Bob",
        "age": 25,
        "hobbies": ["gaming", "cooking"]
    }`

    err = json.Unmarshal([]byte(datastr), &data)
    if err != nil {
        log.Error("Cannot connect to customer service", err)
        api.InternalErrorHandler(w)
        return
    }

	user, err := services.NewCustomer(user_id)
	if err != nil {
        log.Error("Cannot connect to customer service", err)
        api.InternalErrorHandler(w)
        return
	}

	res, err := user.RunProduct(user_id, tool_id, data)
    if err != nil {
        log.Error("Cannot connect to customer service", err)
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

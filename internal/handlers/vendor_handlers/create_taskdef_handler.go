package vendor_handlers
// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
import (
	//"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fizlip/microstore/api"
	log "github.com/sirupsen/logrus"
)

// CreateTaskdef will register a new Task Definition in aws ecs.
func CreateTaskdef(w http.ResponseWriter, r *http.Request) {
	var err error
	var taskdef api.Taskdef

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &taskdef)
	if err != nil {
		panic(err)
	}
	log.Println(taskdef)

	response := api.Success{
		Code:    200,
		Message: "Taskdef created",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

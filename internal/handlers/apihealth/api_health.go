// package f4rmer_handlers provides functions used to interact with f4rmers.
package apihealth

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// AddTool will add a tool to a f4rmers toolbox. Each tool is represented
// by a tool id (tid). To add a tool you need to define the user adding the 
// tool (usr), the f4rmer that should have the tool (uid) and the tool to add (tid).
//
// example: POST https://<hostname>/user/f4rmer/add/tool?usr=<username>&uid=<f4rmer_id>&tid=<tool_id>
// response: 
//      SUCCESS: 200, "tool was added", 
//      FAILURE: 403, An unexpected error occured 
//
func IsLive(w http.ResponseWriter, r *http.Request) {
    var err error

    err = json.NewEncoder(w).Encode("f4rmhouse api is live!")
    if err != nil {
        log.Error(err)
        api.InternalErrorHandler(w)
        return
    }
}

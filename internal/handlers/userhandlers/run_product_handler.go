package userhandlers

import (
	"encoding/json"
	//"github.com/google/uuid"
	"net/http"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	log "github.com/sirupsen/logrus"
)

// RunProduct will execute a tool using the execution engine and update the
// usage of the tool that has been used
func RunProduct(w http.ResponseWriter, r *http.Request) {
	var err error

	user_id := r.Header.Get("X-Username")
	tool_id := r.URL.Query().Get("tid")

	user, err := services.NewCustomer(user_id)
	if err != nil {
		log.Error("Cannot connect to customer service", err)
		api.InternalErrorHandler(w)
		return
	}

	err = user.RunProduct(user_id, tool_id)
	if err != nil {
		log.Error("Cannot connect to customer service", err)
		api.CustomInternalErrorHandler(w, err)
		return
	}

	err = json.NewEncoder(w).Encode("Product usage updated")
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

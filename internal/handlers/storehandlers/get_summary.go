package storehandlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetSummary(w http.ResponseWriter, r *http.Request) {
	var err error

    var tns *services.ToolNameService
    tns, err = services.NewToolNameService()
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

    uti := r.URL.Query().Get("uti")

	result, err := (*tns).GetTool(uti)
    if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

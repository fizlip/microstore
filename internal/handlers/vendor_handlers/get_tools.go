package vendor_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetCurrent...
// ex: <hostname>/user/review/create
// data: ...
func GetTools(w http.ResponseWriter, r *http.Request) {
    vendor, err := services.NewVendorController("fillatino@gmail.com")
	if err != nil {
		log.Fatalf("failed to retrieve usage data, %v", err)
        return
	}

    result, err := vendor.GetTools()
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

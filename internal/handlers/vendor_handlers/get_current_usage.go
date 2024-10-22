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
func GetCurrentUsage(w http.ResponseWriter, r *http.Request) {
	var err error
	vendor, err := services.NewVendorController("example@mail.com")

	if err != nil {
		log.Fatalf("failed to retrieve usage data, %v", err)
	}

	rs, err := vendor.GetToolUsage()

	err = json.NewEncoder(w).Encode(rs)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

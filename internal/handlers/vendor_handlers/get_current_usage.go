// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
package vendor_handlers

import (
	"encoding/json"
	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// GetCurrentUsage returns a struct with the usge of each tool that the vendor has
// created 
//
// example: https://<hostname>/get/usage
// response: 
//      SUCCESS: 200, <usage data>
//      FAILURE: 500, An unexpected error occurred
func GetCurrentUsage(w http.ResponseWriter, r *http.Request) {
	var err error
	vendor, err := services.NewVendorController("fillatino@gmail.com")

	if err != nil {
		log.Fatalf("failed to retrieve usage data, %v", err)
	}

	tool_usage, err := vendor.GetToolUsage()

	err = json.NewEncoder(w).Encode(tool_usage)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

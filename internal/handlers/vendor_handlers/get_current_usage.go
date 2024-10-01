package vendor_handlers 

import (
  "encoding/json"
  "net/http"
  "github.com/fizlip/microstore/api"

  log "github.com/sirupsen/logrus"
)

// GetCurrent...
// ex: <hostname>/user/review/create
// data: ...
func GetCurrentUsage(w http.ResponseWriter, r *http.Request) {
  var err error
  var vendor *api.Vendor
  vendor, err = api.NewVendor("example@mail.com")

  if err != nil {
    log.Fatalf("failed to retrieve usage data, %v", err)
  }

  rs, err := (*vendor).GetProductUsage()

  err = json.NewEncoder(w).Encode(rs)
  if err != nil {
    log.Error(err)
    api.InternalErrorHandler(w)
    return
  }
}

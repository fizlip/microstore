// package vendor_handlers contains request handlers for all endpoints that
// interact with the vendor
package vendorhandlers

import (
    "net/http"
    "encoding/json"

    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// UploadShowcaseImage will upload an image to an S3 bucket.
//
// example: POST https://<hostname>/vendor/image/upload
// response:
//      SUCCESS: 200, Image uploaded successfully
//      FAILURE: 500, An unexepcted error occurred
func UploadShowcaseImage(w http.ResponseWriter, r *http.Request) {
    var err error
    var vendor *services.Vendor

    vendor, err = services.NewVendorController("fillatino@gmail.com")

    err = vendor.UploadShowcaseImage()
    
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }
    
    response := api.Success{
        Code: 200, 
        Message: "Successfully uploaded image",
    }

    err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

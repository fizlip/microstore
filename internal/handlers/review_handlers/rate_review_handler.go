package review_handlers

import (
	"encoding/json"
    "net/http"
    "io/ioutil"

	"github.com/fizlip/microstore/api"
	"github.com/fizlip/microstore/internal/services"

	log "github.com/sirupsen/logrus"
)

// RateReview will ...
// ex: <hostname>/user/review/create
// data: ...
func RateReview(w http.ResponseWriter, r *http.Request) {
	var response api.Success
	var err error
	var reviews *services.ReviewService
    var params api.UpdateReviewParams

    var username = r.Header.Get("X-Username")

	reviews, err = services.NewReviewService()
	if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
        return 
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
        return 
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
	} 

	err = (*reviews).UpdateRating(params.TID, params.UID, username)
    if err != nil {
        log.Error(err)
		api.InternalErrorHandler(w)
		return
    }
	response = api.Success{
		Code:    200,
		Message: "Sucessfully rated review",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

}

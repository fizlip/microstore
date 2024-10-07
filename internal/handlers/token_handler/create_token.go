package token_handler

import (
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/models"
    "io/ioutil"
    "net/http"
    log "github.com/sirupsen/logrus"
    "encoding/json"
)

func InsertToken(w http.ResponseWriter, r *http.Request) {
    var err error
    var token models.Token

    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

    log.Println(token)

    response := api.Success{
		Code:    200,
		Message: "New token created",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}
}

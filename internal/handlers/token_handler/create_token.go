package token_handler

import (
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
    "io/ioutil"
    "net/http"
    log "github.com/sirupsen/logrus"
    "encoding/json"
)

func InsertToken(w http.ResponseWriter, r *http.Request) {
    var err error
    var token services.Token

    body, err := ioutil.ReadAll(r.Body)

    if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}
    log.Println("Origin: ", r.Header.Get("Origin"))
    err = token.Insert()
    if err != nil {
        panic(err)
    }
    if r.Header.Get("Origin") != "http://localhost:3000" {
        response := api.Success{
            Code:    400,
            Message: "User not authorized",
        }

        err = json.NewEncoder(w).Encode(response)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
    } else {
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
}

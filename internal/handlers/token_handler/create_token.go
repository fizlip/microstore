package token_handler

import (
    "github.com/fizlip/microstore/api"
    "github.com/fizlip/microstore/internal/services"
    "io/ioutil"
    "net/http"
    log "github.com/sirupsen/logrus"
    "encoding/json"
    "github.com/joho/godotenv"
    "os"
)

func InsertToken(w http.ResponseWriter, r *http.Request) {
    var err error
    var token services.Token

    err = godotenv.Load(".env")
    if err != nil {
        panic(err)
    }

    // Read body which contains oauth token to be inserted
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
		panic(err)
	}
    // Get token
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

    log.Println("Authorizartion: ", r.Header.Get("Authorization"))

    // Check if secret_key is correct
    if r.Header.Get("Authorization") != os.Getenv("SECRET_KEY") {
        response := api.Success{
            Code:    400,
            Message: "User not authorized",
        }
        log.Error("Failed authorization attempt")
        err = json.NewEncoder(w).Encode(response)
        if err != nil {
            log.Error(err)
            api.InternalErrorHandler(w)
            return
        }
    } else {
        err = token.Insert()
        if err != nil {
            panic(err)
        }
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

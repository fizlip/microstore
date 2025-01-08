// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services 

import (

    "fmt"
    "net/http"
    "time"
    "encoding/json"

    "github.com/fizlip/microstore/internal/models"

    log "github.com/sirupsen/logrus"
)

// RemoteTokenData is used to access tokens from OAuth providers.
// The structs and its functions provide a simple way to get the 
// raw token from OAuth providers 
type RemoteTokenData struct {
    url string 
    header string
    headerValue string
    token string
    requestInstigator string
}

// fetch will send an http request to an OAuth provider and return
// the session object, stored as json in the http.Response struct
func (v *RemoteTokenData) fetch() (*http.Response, error){
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    req, _ := http.NewRequest("GET", v.url, nil)
    if len(v.header) > 0{
        req.Header.Set(v.header, v.headerValue)
    }
    res, err := client.Do(req)
    if err != nil {
        return res, fmt.Errorf("failed to send request: %v", err)
    }

    if res.StatusCode != http.StatusOK {
        return res, fmt.Errorf("token verification failed with status: %d", res.StatusCode)
    }
    return res, nil
}

// validateAccessToken ensures 
// 1. The given token exists 
// 2. The owner of the token is the same as the instigator 
func validateAccessToken(remote RemoteTokenData, unwrapperFunc func(resp *http.Response) (string, error)) (string, error) {
    resp, err := remote.fetch() 
    defer resp.Body.Close()
    // If token does not exists unwrapper will return an err
    email, err := unwrapperFunc(resp)
    if err != nil {
        log.Error("Token does not exist or is invalid")
        return "", fmt.Errorf("Token does not exist")
    }
    if email != remote.requestInstigator {
        return "", fmt.Errorf("The instigator of the request does not match the token owner.")
    }
    return email, err
}

// unwrapGoogleResponse returns the userid from the response
// returned by Google OAuth provider
func unwrapGoogleResponse(resp *http.Response) (string, error) {
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("Error decoding Google response JSON: %v", err)
    }
    if _, ok := result["error"]; ok {
        log.Error("Error response from remote server: ", result["error"])
        return "", fmt.Errorf("Error in remote server response")
    }
    return result["email"].(string), nil
}

// unwrapGithubResponse returns userid from the response 
// returned by GitHub OAuth provider
func unwrapGithubResponse(resp *http.Response) (string, error) {
    var result []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("Error decoding Github response JSON: %v", err)
    }
    return result[0]["email"].(string), nil
}

func ValidateAPIKey(username, key string) (string, error){
    apiKeyT, err := models.NewDynamoTable[models.APIKey]("f4-api-keys", "username", "key")
    if err != nil {
        return "", fmt.Errorf("Error connecting to db: ", err)
    }
    storedKeys, err := apiKeyT.Read(username)
    if err != nil {
        return "", fmt.Errorf("Error getting api key", err)
    }
    for _, el := range storedKeys {
        ok := models.ValidateKey(key, el.Key); 
        if ok {
            return username, nil
        } 
    }
    return "", fmt.Errorf("API key is not valid")
}

// ValidateAccesstoken will return an error is a token is invalid or nil otherwise
func ValidateAccessToken(token string, provider string, instigator string) (string, error) {
    // Currently Github & Google are supported as OAuth providers
    // each OAuth has different ways to validate the token
    githubTokenData := RemoteTokenData{
        "https://api.github.com/user/emails", 
        "Authorization",
        fmt.Sprintf("token %s", token),
        token,
        instigator,
    }

    googleTokenData := RemoteTokenData{
        fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", token),
        "",
        "",
        token,
        instigator,
    }

    // Determine strategy based on OAuth provider
    switch {
        case provider == "google":
            return validateAccessToken(googleTokenData, unwrapGoogleResponse)
        case provider == "github":
            return validateAccessToken(githubTokenData, unwrapGithubResponse)
        case provider == "f4rmhouse":
            return ValidateAPIKey(instigator, token)
        default:
            return "", fmt.Errorf("%s is not recognized as a provider", provider)
    }
}

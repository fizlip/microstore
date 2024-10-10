package services 

import (
    "fmt"
    "net/http"
    "time"
    "encoding/json"
    "io"
)

type TokenValidator interface {
    getTokenInfo(token string) (*http.Response, error)
    getTokenOwner(token string) (string, error)
    validateAccessToken(token string) error
}


type RemoteTokenData struct {
    url string 
    header string
    headerValue string
    token string
    requestInstigator string
}

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

func validateGoogleAccessToken(remote RemoteTokenData) error {
    resp, err := remote.fetch()
    defer resp.Body.Close()

    var result map[string]interface{}
    if errr := json.NewDecoder(resp.Body).Decode(&result); errr != nil {
        fmt.Errorf("Error decoding JSON: %v", errr)
    }
    email, _ := result["email"].(string)
    if email != remote.requestInstigator {
        return fmt.Errorf("The instigator of the request does not match the token owner.")
    }

    return err
}

func validateGithubAccessToken(remote RemoteTokenData) error {
    resp, err := remote.fetch() 
    defer resp.Body.Close()
    r, _ := io.ReadAll(resp.Body) 
    fmt.Println(string(r))
    return err
}

func ValidateAccessToken(token string, provider string, instigator string) error {

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

    switch {
        case provider == "google":
            return validateGoogleAccessToken(googleTokenData)
        case provider == "github":
            return validateGithubAccessToken(githubTokenData)
        default:
            return fmt.Errorf("%s is not recognized as a provider", provider)
    }
}

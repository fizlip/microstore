package services 

import (
    "fmt"
    "net/http"
    "time"
    "encoding/json"
    "io"
)

func getGoogleTokenInfo(token string) (*http.Response, error) {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create request to Google's tokeninfo endpoint
    url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", token)
    resp, err := client.Get(url)
    if err != nil {
        return resp, fmt.Errorf("failed to send request: %v", err)
    }

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return resp, fmt.Errorf("token verification failed with status: %d", resp.StatusCode)
    }
    return resp, nil
}

func getGoogleTokenOwner(token string) (string, error) {
    resp, err := getGoogleTokenInfo(token)
    defer resp.Body.Close()

    var result map[string]interface{}
    if errr := json.NewDecoder(resp.Body).Decode(&result); errr != nil {
        fmt.Errorf("Error decoding JSON: %v", errr)
    }
    email, _ := result["email"].(string)

    return email, err 
}

func validateGoogleAccessToken(token string) error {
    _, err := getGoogleTokenOwner(token)
    return err 
}

func getGithubTokenInfo(token string) (*http.Response, error) {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create request to Google's tokeninfo endpoint
    url := fmt.Sprintf("https://api.github.com/user/emails")
    
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
    res, err := client.Do(req)
    if err != nil {
        return res, fmt.Errorf("failed to send request: %v", err)
    }

    // Check response status
    if res.StatusCode != http.StatusOK {
        return res, fmt.Errorf("token verification failed with status: %d", res.StatusCode)
    }
    return res, nil
}

func getGithubTokenOwner(token string) error {
    resp, err := getGithubTokenInfo(token) 
    defer resp.Body.Close()
    r, _ := io.ReadAll(resp.Body) 
    fmt.Println(string(r))
    return err 
}

func validateGithubAccessToken(token string) error {
    err := getGithubTokenOwner(token)
    fmt.Println(err)
    return err 
}

func ValidateAccessToken(token string, provider string) error {
    switch {
        case provider == "google":
            return validateGoogleAccessToken(token)
        case provider == "github":
            return validateGithubAccessToken(token)
        default:
            return fmt.Errorf("%s is not recognized as a provider", provider)
    }
}

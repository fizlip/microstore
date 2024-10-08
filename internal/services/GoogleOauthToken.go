package services 

import (
    "fmt"
    "net/http"
    "time"
)

func validateGoogleAccessToken(token string) error {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create request to Google's tokeninfo endpoint
    url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", token)
    resp, err := client.Get(url)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("token verification failed with status: %d", resp.StatusCode)
    }
    return nil
}

func validateGithubAccessToken(token string) error {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Create request to Google's tokeninfo endpoint
    url := fmt.Sprintf("https://api.github.com/user/repos")
    
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
    res, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }
    defer res.Body.Close()

    // Check response status
    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("token verification failed with status: %d", res.StatusCode)
    }
    return nil
}

func ValidateAccessToken(token string, provider string) error {
    fmt.Println("Provider", provider)
    switch {
        case provider == "google":
            return validateGoogleAccessToken(token)
        case provider == "github":
            return validateGithubAccessToken(token)
        default:
            return fmt.Errorf("%s is not recognized as a provider", provider)
    }
}

package services

import (
    "context"
    "fmt"
    "os"
    "strings"
    "encoding/base64"
    "encoding/json"
    "io"
    "errors"

    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types/image"
    "github.com/docker/docker/api/types/registry"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ecr"

    log "github.com/sirupsen/logrus"
)

func PushToECR(dockerClient *client.Client, imageName string) error {
    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        return fmt.Errorf("unable to load SDK config: %v", err)
    }

    // Create ECR client
    client := ecr.NewFromConfig(cfg)

    // Get ECR authentication token
    input := &ecr.GetAuthorizationTokenInput{}
    result, err := client.GetAuthorizationToken(context.TODO(), input)
    if err != nil {
        return fmt.Errorf("unable to get auth token: %v", err)
    }

    // Extract auth data
    authData := result.AuthorizationData[0]
    registryURL := *authData.ProxyEndpoint
    authToken := *authData.AuthorizationToken

    // Authenticate Docker with ECR
    authStr, err := dockerLogin(dockerClient, registryURL, authToken)
    if err != nil {
        return err
    }

    // Tag the Docker image
    tag, err := dockerTag(dockerClient, imageName, registryURL, "")
    if err != nil {
        log.Error("Error creating tag")
        return err
    }

    // Push the image to ECR
    if err := dockerPush(dockerClient, authStr, tag, registryURL); err != nil {
        return err
    }

    return nil
}

func DockerLoad(dockerCli *client.Client, pathToImage string) error {
    imageTar, err := os.Open(pathToImage)
    if err != nil {
        return err
    }

    defer imageTar.Close()

    // Load the image
    ctx := context.Background()
    loadResponse, err := dockerCli.ImageLoad(ctx, imageTar, true)
    if err != nil {
        return err
    }

    defer loadResponse.Body.Close()

    log.Println("Loaded docker image ", pathToImage)
    log.Println("Loaded respoonse ", loadResponse)

    return nil 
}

func DockerList(dockerCli *client.Client) error {
    images, err := dockerCli.ImageList(context.Background(), image.ListOptions{}) 
    if err != nil {
        return err
    }

    for _, image := range images {
        log.Printf("ID: %s\n", image.ID)
        if len(image.RepoTags) > 0 {
            log.Printf("Tags: %v\n", image.RepoTags)
        }
        log.Printf("Size: %d bytes\n", image.Size)
        log.Printf("Created: %d\n\n", image.Created)
    }

    return nil
}

func dockerLogin(dockerClient *client.Client, registryURL, authToken string) (string, error) {

    decodedToken, err := base64.StdEncoding.DecodeString(authToken)
    if err != nil {
        return "", err
    }

    parts := strings.SplitN(string(decodedToken), ":", 2)
    username := parts[0]
    password := parts[1]

    auth := registry.AuthConfig{
        Username:      username,
        Password:      password,
        ServerAddress: registryURL,
    }

    _, err = dockerClient.RegistryLogin(context.Background(), auth)
    if err != nil {
        return "", err
    }

    encodedJSON, err := json.Marshal(auth)
    if err != nil {
        return "", err
    }

    authStr := base64.URLEncoding.EncodeToString(encodedJSON)

    log.Println("Login succeeded")
    return authStr, nil
}

func dockerTag(dockerClient *client.Client, imageName, registryURL, tag string) (string, error) {
    targetImage := fmt.Sprintf("%s/f4rmhouse:%s", strings.TrimPrefix(registryURL, "https://"), imageName)
    
    if err := dockerClient.ImageTag(context.Background(), imageName, targetImage); err != nil {
        return "", err
    }

    log.Println("Docker tag suceeded. Target: ", targetImage)
    return targetImage, nil
}

func dockerPush(dockerClient *client.Client, authString, targetImage, registryURL string) error {
    pushOpts := image.PushOptions{
        All: true,
        RegistryAuth: authString,
    }
    pushResponse, err := dockerClient.ImagePush(context.Background(), targetImage, pushOpts)
    if err != nil {
        return err
    }
    defer pushResponse.Close()
    if err = processResponse(pushResponse); err != nil {
        log.Error("Failed to push image")
        return err
    }

    log.Println("Docker push suceeded", targetImage)
    return nil
}

func processResponse(r io.Reader) error {
    decoder := json.NewDecoder(r)
    for {
        var message map[string]interface{}
        if err := decoder.Decode(&message); err != nil {
            if err == io.EOF {
               break 
            }
            return err
        }
        if _, ok := message["error"]; ok {
            return errors.New("Error when pushing")
        }
        if status, ok := message["statys"]; ok {
            log.Println(status)
        }
    }
    return nil
}

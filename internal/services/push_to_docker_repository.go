package services

import (
    "context"
    "fmt"
    "os/exec"
    "os"
    "strings"
    "encoding/base64"

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
    if err := dockerLogin(dockerClient, registryURL, authToken); err != nil {
        return err
    }

    //// Tag the Docker image
    if err := dockerTag(dockerClient, imageName, registryURL, "version1.0"); err != nil {
        return err
    }

    //// Push the image to ECR
    //if err := dockerPush(imageName, registryURL); err != nil {
    //    return err
    //}

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

func dockerLogin(dockerClient *client.Client, registryURL, authToken string) error {

    decodedToken, err := base64.StdEncoding.DecodeString(authToken)
    if err != nil {
        return err
    }

    parts := strings.SplitN(string(decodedToken), ":", 2)
    username := parts[0]
    password := parts[1]

    auth := registry.AuthConfig{
        Username:      username,
        Password:      password,
        ServerAddress: registryURL,
    }

    authString, err := dockerClient.RegistryLogin(context.Background(), auth)
    if err != nil {
        return err
    }
    log.Println("Returned from RegistryLogin: ", authString)
    return nil
}

func dockerTag(dockerClient *client.Client, imageName, registryURL, tag string) error {
    targetImage := fmt.Sprintf("%s/%s:%s", strings.TrimPrefix(registryURL, "https://"), imageName, tag)
    
    if err := dockerClient.ImageTag(context.Background(), imageName, targetImage); err != nil {
        return err
    }

    log.Println("Docker tag suceeded. Target: ", targetImage)
    return nil
}

func dockerPush(imageName, registryURL string) error {
    cmd := exec.Command("docker", "push", fmt.Sprintf("%s/%s", registryURL, imageName))
    return cmd.Run()
}

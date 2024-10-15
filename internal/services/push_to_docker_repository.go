package services

import (
    "context"
    "fmt"
    "os/exec"
    "os"
    "strings"

    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types/image"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ecr"

    log "github.com/sirupsen/logrus"
)

func PushToECR(imageName string) error {
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

    log.Println("authData: ", authData)
    log.Println("registryURL: ", registryURL)
    log.Println("authToken: ", authToken)

    // Authenticate Docker with ECR
    //if err := dockerLogin(registryURL, authToken); err != nil {
    //    return err
    //}

    //// Tag the Docker image
    //if err := dockerTag(imageName, registryURL); err != nil {
    //    return err
    //}

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

func dockerLogin(registryURL, authToken string) error {
    cmd := exec.Command("docker", "login", "-u", "AWS", "--password-stdin", registryURL)
    cmd.Stdin = strings.NewReader(authToken)
    return cmd.Run()
}

func dockerTag(imageName, registryURL string) error {
    cmd := exec.Command("docker", "tag", imageName, fmt.Sprintf("%s/%s", registryURL, imageName))
    return cmd.Run()
}

func dockerPush(imageName, registryURL string) error {
    cmd := exec.Command("docker", "push", fmt.Sprintf("%s/%s", registryURL, imageName))
    return cmd.Run()
}

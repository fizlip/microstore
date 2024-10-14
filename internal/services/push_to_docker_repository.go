package services

import (
    "context"
    "fmt"
    "os/exec"
    "strings"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ecr"
)

func pushToECR(imageName string) error {
    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
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
    if err := dockerLogin(registryURL, authToken); err != nil {
        return err
    }

    // Tag the Docker image
    if err := dockerTag(imageName, registryURL); err != nil {
        return err
    }

    // Push the image to ECR
    if err := dockerPush(imageName, registryURL); err != nil {
        return err
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

package services

import (
    "fmt"
    "context"
    "errors"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ecs"
    "github.com/aws/aws-sdk-go-v2/service/ecs/types"

    log "github.com/sirupsen/logrus"
)

type DeployService struct {
  cfg   aws.Config 
}

func (d *DeployService) CreateCluster(clusterName string) (error) {
    log.Println("Creating cluster...")    

    // Create an ECS client
    client := ecs.NewFromConfig(d.cfg)

    input := &ecs.CreateClusterInput{
        ClusterName: &clusterName,
        
        // Optional: Add tags to the cluster
        Tags: []types.Tag{
            {
                Key:   aws.String("Environment"),
                Value: aws.String("Production"),
            },
        },

        // Optional: Configure cluster settings
        Settings: []types.ClusterSetting{
            {
                Name:  types.ClusterSettingNameContainerInsights,
                Value: aws.String("enabled"),
            },
        },

        // Optional: Specify capacity providers
        CapacityProviders: []string{
            "FARGATE",
            "FARGATE_SPOT",
        },
    }

    // Create the ECS cluster
    _, err := client.CreateCluster(context.TODO(), input)
    if err != nil {
        return errors.New(fmt.Sprintf("failed to create ECS cluster: %v", err))
    }

    log.Println("Cluster has been created")

    return err 
}

func (d *DeployService) CreateTaskDef() (error) {
    log.Println("Creating TaskDef...")    
    return nil
}

func (d *DeployService) CreateService() (error) {
    log.Println("Creating Service...")    
    return nil
}

func NewDeployService() (*DeployService, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        return &DeployService{}, err
    } 
    return &DeployService{cfg: cfg}, nil
}

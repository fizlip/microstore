package services

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"

    log "github.com/sirupsen/logrus"
)

type DeployService struct {
  cfg   aws.Config 
}

func (d *DeployService) CreateCluster(clusterName string) (error) {
    log.Println("Creating cluster...")    
    return nil
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

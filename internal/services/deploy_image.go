package services

import (
    "fmt"
    "context"
    "errors"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ecs"
    "github.com/aws/aws-sdk-go-v2/service/ecs/types"

    log "github.com/sirupsen/logrus"
)

type DeployService struct {
  cfg           aws.Config 
  clusterName   string
}

func (d *DeployService) GetEndpoint(clusterName string) ([]string, error) {
    log.Println("Get endpoint")    

    var publicIP []string
    // Create ECS client
    ecsClient := ecs.NewFromConfig(d.cfg)
    ec2Client := ec2.NewFromConfig(d.cfg)

    // List all tasks in the cluster
    listTasksInput := &ecs.ListTasksInput{
        Cluster: &clusterName,
    }
    tasksResp, err := ecsClient.ListTasks(context.TODO(), listTasksInput)
    if err != nil {
        return []string{}, errors.New(fmt.Sprintf("error listing tasks: %v", err))
    }

    // If no tasks found
    if len(tasksResp.TaskArns) == 0 {
        log.Printf("No tasks found in cluster %s\n", clusterName)
        return nil, errors.New(fmt.Sprintf("No tasks found in cluster %s\n", clusterName))
    }

    // Describe the tasks to get more details
    describeTasksInput := &ecs.DescribeTasksInput{
        Cluster: &clusterName,
        Tasks:   tasksResp.TaskArns,
    }
    tasksDesc, err := ecsClient.DescribeTasks(context.TODO(), describeTasksInput)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("error describing tasks: %v", err))
    }

    // Print task details with network information
    log.Printf("Tasks in cluster %s:\n", clusterName)
    for _, task := range tasksDesc.Tasks {
        // Print container details
        for _, container := range task.Containers {
            if string(*container.LastStatus) != "RUNNING" {
                publicIP = append(publicIP, string(*container.LastStatus))
                return publicIP, nil
            }
        }
        
        // Retrieve network interfaces to get public IP
        if task.Attachments != nil {
            log.Println("Attachements: ")
            for _, attachment := range task.Attachments {
                log.Println(*attachment.Type)
                if *attachment.Type == "ElasticNetworkInterface" {
                    // Extract the network interface ID from attachments
                    var networkInterfaceId string
                    for _, detail := range attachment.Details {
                        if *detail.Name == "networkInterfaceId" {
                            networkInterfaceId = *detail.Value
                            break
                        }
                    }
                    
                    // Describe network interface to get public IP
                    if networkInterfaceId != "" {
                        networkInterfacesInput := &ec2.DescribeNetworkInterfacesInput{
                            NetworkInterfaceIds: []string{networkInterfaceId},
                        }
                        
                        networkInterfacesResp, err := ec2Client.DescribeNetworkInterfaces(context.TODO(), networkInterfacesInput)
                        if err != nil {
                            log.Printf("    Error retrieving network interface: %v\n", err)
                            continue
                        }
                        
                        // Print public IP if available
                        if len(networkInterfacesResp.NetworkInterfaces) > 0 {
                            ni := networkInterfacesResp.NetworkInterfaces[0]
                            if ni.Association != nil && ni.Association.PublicIp != nil {
                                log.Printf("    Public IP: %s\n", *ni.Association.PublicIp)
                                publicIP = append(publicIP, *ni.Association.PublicIp)
                            }
                        }
                    }
                }
            }
        }
        
        log.Println("---")
    }

    return publicIP, nil
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
    d.clusterName = clusterName

    return err 
}

func (d *DeployService) CreateTaskDef(taskDef, tool_name, tool_uri string) (error) {
    log.Println("Creating TaskDef...")    

    client := ecs.NewFromConfig(d.cfg)

    containerDefinitions := []types.ContainerDefinition{
        {
            Name:  aws.String(tool_name),
            Image: aws.String(tool_uri),
            
            // CPU and memory settings
            Cpu:    *aws.Int32(256),  // 0.25 vCPU
            Memory: aws.Int32(512),  // 512 MiB

            // Port mappings (optional)
            PortMappings: []types.PortMapping{
                {
                    ContainerPort: aws.Int32(80),
                    HostPort:      aws.Int32(80),
                    Protocol:      types.TransportProtocolTcp,
                },
                {
                    ContainerPort: aws.Int32(8080),
                    HostPort:      aws.Int32(8080),
                    Protocol:      types.TransportProtocolTcp,
                },
            },

            // Environment variables (optional)
            Environment: []types.KeyValuePair{
                {
                    Name:  aws.String("ENV"),
                    Value: aws.String("production"),
                },
            },

            // Essential flag - if true, the entire task fails if this container fails
            Essential: aws.Bool(true),
        },
    }

    // Prepare task definition input
    input := &ecs.RegisterTaskDefinitionInput{
        Family: aws.String(taskDef),
        
        // Container definitions
        ContainerDefinitions: containerDefinitions,
        
        // Task execution role (for pulling ECR images, sending logs, etc.)
        ExecutionRoleArn: aws.String("arn:aws:iam::236798077409:role/ecsTaskExecutionRole"),
        
        // Task role (for application permissions)
        TaskRoleArn: aws.String("arn:aws:iam::236798077409:role/ecsTaskExecutionRole"),
        
        // Network mode
        NetworkMode: types.NetworkModeAwsvpc,
        
        // CPU and memory at the task level
        Cpu:    aws.String("256"),
        Memory: aws.String("512"),

        RuntimePlatform: &types.RuntimePlatform{
            OperatingSystemFamily: types.OSFamilyLinux,
            CpuArchitecture:       types.CPUArchitectureX8664,
        },
        
        // Optional: Add tags
        Tags: []types.Tag{
            {
                Key:   aws.String("Environment"),
                Value: aws.String("Production"),
            },
        },

        // Optional: Requires compatibility
        RequiresCompatibilities: []types.Compatibility{
            types.CompatibilityFargate,
        },
    }

    _, err := client.RegisterTaskDefinition(context.TODO(), input)

    return err 
}

func (d *DeployService) CreateService(clusterName, serviceName, taskDefinition string) (error) {
    log.Println("Creating Service...")    

    client := ecs.NewFromConfig(d.cfg)

    serviceInput := &ecs.CreateServiceInput{
        // Cluster and service basic configuration
        Cluster:        aws.String(clusterName),
        ServiceName:    aws.String(serviceName),
        TaskDefinition: aws.String(taskDefinition),

        // Deployment configuration
        DesiredCount: aws.Int32(2), // Number of task instances to run
        LaunchType:   types.LaunchTypeFargate,

        // Network configuration (required for Fargate)
        NetworkConfiguration: &types.NetworkConfiguration{
            AwsvpcConfiguration: &types.AwsVpcConfiguration{
                Subnets: []string{
                    "subnet-fd1bbd97",
                    "subnet-7d917231",
                    "subnet-8e001cf3",
                },
                SecurityGroups: []string{
                    "sg-3f2ecd5a", // Replace with your security group IDs
                    "sg-02e04cfdc02886613",
                },
                AssignPublicIp: types.AssignPublicIpEnabled,
            },
        },

        // Optional: Deployment strategy
        DeploymentConfiguration: &types.DeploymentConfiguration{
            MaximumPercent:        aws.Int32(200),
            MinimumHealthyPercent: aws.Int32(50),
        },

        // Optional: Service tags
        Tags: []types.Tag{
            {
                Key:   aws.String("Environment"),
                Value: aws.String("Production"),
            },
        },
    }

    _, err := client.CreateService(context.TODO(), serviceInput)

    return err 
}

func NewDeployService() (*DeployService, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
    if err != nil {
        return &DeployService{}, err
    } 
    return &DeployService{cfg: cfg}, nil
}

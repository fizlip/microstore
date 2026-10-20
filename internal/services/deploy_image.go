// package services contains implementaions of models and interfaces
// and provides the business logic of the API
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	log "github.com/sirupsen/logrus"
)

// DeployService is used to deploy a service to ecs. This means that
// hardware will be allocated to the clust that has been named
type DeployService struct {
	cfg         aws.Config
	clusterName string
}

// GetEndpoint will return all IP addresses associated with this clusterName
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

// CreateCluster will create a new ECS cluster with this clusterName
func (d *DeployService) CreateCluster(clusterName string) error {
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

// CreateTaskDef creates a ECS task definition. A Task Definition is a bluebrint
// given to the ECS platform that tells it:
//
// * What docker container to use
// * How much CPU and memory to allocate
// * Which ports to expose
// * How containers should be linked together
// * What environment variables to set
// * Which storage volumes to mount
// * The networking mode to use
//
// Think of it as a template that defines everything needed to run your application.
// A task definition can specify one or more containers that work together. For example,
// you might have a task definition that includes both a web application container
// and a database container.
func (d *DeployService) CreateTaskDef(taskDef, tool_name, tool_uri string, secrets []types.KeyValuePair) error {

	client := ecs.NewFromConfig(d.cfg)

	containerDefinitions := []types.ContainerDefinition{
		{
			Name:  aws.String(tool_name),
			Image: aws.String(tool_uri),

			// CPU and memory settings
			Cpu:    *aws.Int32(256), // 0.25 vCPU
			Memory: aws.Int32(512),  // 512 MiB

			// Port mappings (optional)
			PortMappings: []types.PortMapping{
				{
					ContainerPort: aws.Int32(8080),
					HostPort:      aws.Int32(8080),
					Protocol:      types.TransportProtocolTcp,
				},
			},

			// Environment variables (optional)
			Environment: secrets,

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
		ExecutionRoleArn: aws.String(""),

		// Task role (for application permissions)
		TaskRoleArn: aws.String(""),

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
	if err != nil {
		log.Error(err)
		return errors.New(fmt.Sprintf("error creating taskdef: %v.", err))
	}

	return err
}

func (d *DeployService) CreateService(clusterName, serviceName, taskDefinition string) error {
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
					"",
					"",
					"",
				},
				SecurityGroups: []string{
					"", // Replace with your security group IDs
					"",
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

func (d *DeployService) DeleteCluster(clusterName string) error {
	client := ecs.NewFromConfig(d.cfg)

	listServicesInput := &ecs.ListServicesInput{
		Cluster: &clusterName,
	}

	services, err := client.ListServices(context.TODO(), listServicesInput)
	if err != nil {
		log.Error("Error listing services")
		return errors.New(fmt.Sprintf("error listing services: %v.", err))
	}

	if len(services.ServiceArns) > 0 {
		log.Error("Cluster still has active services")
		return errors.New(fmt.Sprintf("cluster %v still has active services.", err))

	}

	// Delete the cluster
	input := &ecs.DeleteClusterInput{
		Cluster: &clusterName,
	}

	_, err = client.DeleteCluster(context.TODO(), input)
	if err != nil {
		log.Error("Error deleting cluster")
		return errors.New(fmt.Sprintf("error deleting cluster: %v", err))
	}

	log.Println("cluster deleted: ", clusterName)

	return nil
}

func (d *DeployService) DeleteTaskDef(clusterName string) error {

	client := ecs.NewFromConfig(d.cfg)

	taskDefinitions, err := listTaskDefinitions(context.TODO(), client, clusterName+"-taskdef")
	if err != nil {
		return err
	}

	if len(taskDefinitions) == 0 {
		log.Println("No task definitions found")
		return nil
	}

	// Deregister each task definition
	for _, taskDef := range taskDefinitions {
		input := &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: &taskDef,
		}

		_, err := client.DeregisterTaskDefinition(context.TODO(), input)
		if err != nil {
			log.Printf("Warning: Failed to deregister task definition %s: %v\n", taskDef, err)
			continue
		}
		log.Printf("Successfully deregistered task definition: %s\n", taskDef)
	}

	return nil

}

func (d *DeployService) DeleteService(clusterName string) error {
	client := ecs.NewFromConfig(d.cfg)

	services, err := listServices(client, clusterName)
	if err != nil {
		log.Error("Error listing services")
		return errors.New(fmt.Sprintf("error listing services: %v.", err))
	}

	if len(services) == 0 {
		log.Println("No services to delete")
		return nil

	}

	for _, service := range services {
		if err := deleteService(client, clusterName, service); err != nil {
			log.Error("Error deleting service %s: %v", service, err)
		} else {
			log.Printf("Successfully initiated deletion of service: %s", service)
		}
	}

	return nil
}

func NewDeployService() (*DeployService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		return &DeployService{}, err
	}
	return &DeployService{cfg: cfg}, nil
}

func listServices(client *ecs.Client, clusterName string) ([]string, error) {
	var services []string
	var nextToken *string

	for {
		input := &ecs.ListServicesInput{
			Cluster:    &clusterName,
			MaxResults: aws.Int32(100),
			NextToken:  nextToken,
		}

		result, err := client.ListServices(context.TODO(), input)
		if err != nil {
			return nil, fmt.Errorf("failed to list services: %w", err)
		}

		services = append(services, result.ServiceArns...)

		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return services, nil
}

func deleteService(client *ecs.Client, clusterName, serviceArn string) error {
	// First, update the service to have 0 desired count
	_, err := client.UpdateService(context.TODO(), &ecs.UpdateServiceInput{
		Cluster:      &clusterName,
		Service:      &serviceArn,
		DesiredCount: aws.Int32(0),
	})
	if err != nil {
		log.Error("failed to delete service")
		return errors.New(fmt.Sprintf("failed to delete service: %v.", err))
	}

	// Then delete the service
	_, err = client.DeleteService(context.TODO(), &ecs.DeleteServiceInput{
		Cluster: &clusterName,
		Service: &serviceArn,
		Force:   aws.Bool(true), // Use force deletion if needed
	})
	if err != nil {
		log.Error("failed to delete service")
		return errors.New(fmt.Sprintf("failed to delete service: %v.", err))
	}

	// Create a custom waiter for service deletion
	for {
		// Check if service still exists
		output, err := client.DescribeServices(context.TODO(), &ecs.DescribeServicesInput{
			Cluster:  &clusterName,
			Services: []string{serviceArn},
		})

		if err != nil {
			log.Error("error checking service status")
			return errors.New(fmt.Sprintf("error checking service status: %v.", err))
		}

		log.Println("services", output.Services)

		log.Printf("services %v", *output.Services[0].Status)

		// If no services returned or service is inactive, consider it deleted
		if len(output.Services) == 0 ||
			(len(output.Services) == 1 && *output.Services[0].Status == "INACTIVE") {
			break
		}

		// Wait before next check
		time.Sleep(5 * time.Second)
		log.Println("Waiting to finish deletion of service...")
	}

	return nil
}

func listTaskDefinitions(ctx context.Context, client *ecs.Client, clusterName string) ([]string, error) {
	var taskDefinitions []string
	var nextToken *string

	for {
		input := &ecs.ListTaskDefinitionsInput{
			NextToken: nextToken,
			// Optional: Add filters if you want to delete specific families
			FamilyPrefix: aws.String(clusterName),
		}

		output, err := client.ListTaskDefinitions(ctx, input)
		if err != nil {
			log.Error("failed to list task definitions")
			return nil, errors.New(fmt.Sprintf("failed to list task definitions: %v.", err))
		}

		taskDefinitions = append(taskDefinitions, output.TaskDefinitionArns...)

		// Break if no more pages
		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return taskDefinitions, nil
}

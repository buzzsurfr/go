package awsutil

import (
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

// ComputeType is an enum for different compute types
type ComputeType int

const (
	// ComputeUnknown is when the compute type could not be identified
	ComputeUnknown ComputeType = iota
	// ComputeEC2 is running in an EC2 instance
	ComputeEC2
	// ComputeDocker is running in a docker container (and is not ECS, EKS, or kubernetes)
	ComputeDocker
	// ComputeKubernetes is running in a kubernetes pod (and is not EKS)
	ComputeKubernetes
	// ComputeECS is running in an ECS task
	ComputeECS
	// ComputeEKS is running in an EKS pod
	ComputeEKS
	// ComputeLambda is running in a Lambda function
	ComputeLambda
)

// String is the string representation of the ComputeType
func (ct ComputeType) String() string {
	return [...]string{"unknown", "ec2", "docker", "kubernetes", "ecs", "eks", "lambda"}[ct]
}

// DetectOutput is the output struct for the Detect() function
type DetectOutput struct {
	Type ComputeType
}

// String is the string representation of the ComputeType
func (do DetectOutput) String() string {
	return do.Type.String()
}

// IsEC2 returns whether this is running in an EC2 instance
func (do DetectOutput) IsEC2() bool {
	return do.Type == ComputeEC2
}

// IsDocker returns whether this is running in a docker container (and is not ECS, EKS, or kubernetes)
func (do DetectOutput) IsDocker() bool {
	return do.Type == ComputeDocker
}

// IsKubernetes returns whether this is running in a kubernetes pod (and is not EKS)
func (do DetectOutput) IsKubernetes() bool {
	return do.Type == ComputeKubernetes
}

// IsECS returns whether this is running in an ECS task
func (do DetectOutput) IsECS() bool {
	return do.Type == ComputeECS
}

// IsEKS returns whether this is running in an EKS pod
func (do DetectOutput) IsEKS() bool {
	return do.Type == ComputeEKS
}

// IsLambda returns whether this is running in a Lambda function
func (do DetectOutput) IsLambda() bool {
	return do.Type == ComputeLambda
}

// IsEC2 returns whether this is running in an EC2 instance
func IsEC2(do DetectOutput) bool {
	return do.IsEC2()
}

// IsDocker returns whether this is running in a docker container (and is not ECS, EKS, or kubernetes)
func IsDocker(do DetectOutput) bool {
	return do.IsDocker()
}

// IsKubernetes returns whether this is running in a kubernetes pod (and is not EKS)
func IsKubernetes(do DetectOutput) bool {
	return do.IsKubernetes()
}

// IsECS returns whether this is running in an ECS task
func IsECS(do DetectOutput) bool {
	return do.IsECS()
}

// IsEKS returns whether this is running in an EKS pod
func IsEKS(do DetectOutput) bool {
	return do.IsEKS()
}

// IsLambda returns whether this is running in a Lambda function
func IsLambda(do DetectOutput) bool {
	return do.IsLambda()
}

// Detect attempts to determine the compute type of the currently running code.
func Detect() DetectOutput {
	do := DetectOutput{}
	switch {
	case checkLambda():
		do.Type = ComputeLambda
	case checkECS():
		do.Type = ComputeECS
	// case checkEKS():
	// 	do.Type = ComputeEKS
	case checkKubernetes():
		do.Type = ComputeKubernetes
	case checkDocker():
		do.Type = ComputeDocker
	case checkEC2():
		do.Type = ComputeEC2
	default:
		do.Type = ComputeUnknown
	}
	return do
}

// Lambda has two environment variables
func checkLambda() bool {
	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		return true
	}
	return strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda_")
}

func checkECS() bool {
	return strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_ECS_")
}

// Kubernetes sets environment variables and creates a specific folder
func checkKubernetes() bool {
	// KUBERNETES_SERVICE_HOST environment variable
	if _, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		return true
	}
	// folder /var/run/secrets/kubernetes.io exists?
	_, err := os.Stat("/var/run/secrets/kubernetes.io")
	return err != nil && os.IsExist(err)
}

// Docker containers typically contain "docker" in the cgroup names
func checkDocker() bool {
	cgroupBytes, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		return false
	}
	cgroup := string(cgroupBytes)
	return strings.Contains(cgroup, "docker")
}

// checkEC2 inspects the IMDS endpoint (instance metadata) for a valid network connection
func checkEC2() bool {
	conn, err := net.DialTimeout("tcp4", "169.254.169.254:80", 25*time.Millisecond)
	defer conn.Close()
	return err != nil
}

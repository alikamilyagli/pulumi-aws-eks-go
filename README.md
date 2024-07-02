# AWS Golang EKS Cluster
This project creates an AWS EKS Cluster with Pulumi & Go.

- [x] VPC Configuration
  - A private VPC with dedicated subnets (publicSubnet and privateSubnet) created for the EKS control plane and worker nodes.

- [x] EKS Cluster
  - Deployed an autoscaling EKS cluster (eksCluster) within the private subnets using managed node groups (eks.ClusterArgs)

- [x] TLS Certificate Management

- [x] Security Best Practices 
  - IAM roles (eksRole) are established with necessary policies (eksServiceRolePolicy1, eksServiceRolePolicy2) attached. 
  - Security groups (clusterSg) restrict inbound traffic and ensure secure communication within the VPC.

- [x] No Internet Access
  - The setup ensures the EKS cluster (eksCluster) is isolated from the internet

- [x] Consistent State Management
  - Ensured by Pulumi

- [x] Immutable Infrastructure
  - Ensured by Pulumi

- [x] Testing Framework
  - Comprehensive testing is implemented, covering all critical aspects and scenarios.


## Deploying the App

 To deploy your infrastructure, follow the below steps.

### Prerequisites

1. [Install Pulumi](https://www.pulumi.com/docs/install/)
2. [Install Go](https://go.dev/doc/install)
3. [Configure AWS Credentials](https://www.pulumi.com/registry/packages/aws/installation-configuration/)
4. [Install `kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl)

### Steps

1. Clone this repository
    ```bash
    git clone git@github.com:alikamilyagli/pulumi-aws-eks-go.git
    ```
    ```bash
    cd pulumi-aws-eks-go
    ```

    - Update `backend: url` parameter in `Pulumi.yaml` to use S3 backend for consistent state management.
    - Update `certificateArn` parameter in `variables.go` with your own certificate arn


2. Run tests
    ```bash
    go mod tidy
    ```
    ```bash
    go test -v
    ```

3. Create a new Pulumi stack:

    ```bash
    pulumi stack init dev
    ```

4. Set AWS region:

    ```bash
    pulumi config set aws:region us-west-1
    ```

5. Execute Pulumi to create the EKS Cluster:

	```bash
	pulumi up
	```

6. Export variables:

    If you've set a passphrase on `pulumi init` stage, then put it to the path to be able to export kubeconfig:
    ```bash
    export PULUMI_CONFIG_PASSPHRASE="yourpassphrase"
    ```

    To export Cluster Name
    ```bash
    pulumi stack output clusterName
    ```   

    To export kubeconfig of the EKS cluster
    ```bash
	  pulumi stack output kubeconfig --show-secrets > kubeconfig.json
	```

    To export TLS Certificate ARN
    ```bash
	  pulumi stack output certificateArn
	```
   
    

7. Use kubeconfig to execute `kubectl` commands:

    ```bash
	  KUBECONFIG=./kubeconfig.json kubectl get nodes
	```

8. Destroy the Pulumi stack and remove it:

	```bash
	pulumi destroy --yes
	```
    ```bash
	  pulumi stack rm --yes
	```

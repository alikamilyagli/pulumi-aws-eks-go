package main

const (
	region            = "us-west-1"
	availabilityZone1 = "us-west-1a"
	availabilityZone2 = "us-west-1c"
	eksClusterName    = "my-eks-cluster"
	vpcCIDR           = "10.0.0.0/16"
	privateSubnetCIDR = "10.0.1.0/24"
	publicSubnetCIDR  = "10.0.2.0/24"
	clusterNamespace  = "my-namespace"
	instanceType      = "t3.nano"
	domainAddress     = "example.com"
	desiredSize       = 2
	maxSize           = 3
	minSize           = 1

	eksServiceRolePolicy1 = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
	eksServiceRolePolicy2 = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
	nodeGroupPolicy1      = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
	nodeGroupPolicy2      = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
	nodeGroupPolicy3      = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"

	httpsPort = 443
	sshPort   = 22

	vpcName                 = "eks-vpc"
	publicSubnetName        = "eks-public-subnet"
	privateSubnetName       = "eks-private-subnet"
	clusterSgName           = "cluster-sg"
	eksRoleName             = "eks-iam-eksRole"
	nodeGroupRoleName       = "eks-iam-nodeGroupRole"
	attachPolicy1Name       = "eks-iam-attach-policy-1"
	attachPolicy2Name       = "eks-iam-attach-policy-2"
	nodeGroupName           = "eks-node-group"
	eksClusterResourceName  = "eks-cluster"
	nodeGroupPolicyAttachment1  = "nodegroup-policy-attachment-1"
	nodeGroupPolicyAttachment2  = "nodegroup-policy-attachment-2"
	nodeGroupPolicyAttachment3  = "nodegroup-policy-attachment-3"

	EksRoleAssumeRolePolicy = `{
		"Version": "2008-10-17",
		"Statement": [{
			"Sid": "",
			"Effect": "Allow",
			"Principal": {
				"Service": "eks.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}]
	}`
	NodeGroupRoleAssumeRolePolicy = `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {
				"Service": "ec2.amazonaws.com"
			},
			"Action": "sts:AssumeRole"
		}]
	}`
)

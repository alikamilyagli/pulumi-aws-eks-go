package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Main function
func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := SetupInfrastructure(ctx)
        if err != nil {
            return err
        }
        return nil
    })
}

// SetupInfrastructure sets up AWS infrastructure.
func SetupInfrastructure(ctx *pulumi.Context) (*Infrastructure, error) {
	// Create VPC
	vpc, err := ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(vpcCIDR),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name": pulumi.String(vpcName),
		},
	})
	if err != nil {
		return nil, err
	}

    // Create Public Subnet
    publicSubnet, err := ec2.NewSubnet(ctx, publicSubnetName, &ec2.SubnetArgs{
        AvailabilityZone: pulumi.String(availabilityZone1),
        CidrBlock:        pulumi.String(publicSubnetCIDR),
        VpcId:            vpc.ID(),
        Tags: pulumi.StringMap{
            "Name": pulumi.String(publicSubnetName),
        },
    })
    if err != nil {
        return nil, err
    }

    // Create Private Subnet
    privateSubnet, err := ec2.NewSubnet(ctx, privateSubnetName, &ec2.SubnetArgs{
        AvailabilityZone: pulumi.String(availabilityZone2),
        CidrBlock:        pulumi.String(privateSubnetCIDR),
        VpcId:            vpc.ID(),
        Tags: pulumi.StringMap{
            "Name": pulumi.String(privateSubnetName),
        },
    })
    if err != nil {
        return nil, err
    }

    // Create a Security Group for the cluster
    clusterSg, err := ec2.NewSecurityGroup(ctx, clusterSgName, &ec2.SecurityGroupArgs{
        VpcId: vpc.ID(),
        Egress: ec2.SecurityGroupEgressArray{
            ec2.SecurityGroupEgressArgs{
                Protocol:   pulumi.String("-1"),
                FromPort:   pulumi.Int(0),
                ToPort:     pulumi.Int(0),
                CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
            },
        },
        Ingress: ec2.SecurityGroupIngressArray{
            ec2.SecurityGroupIngressArgs{
                Protocol:   pulumi.String("tcp"),
                FromPort:   pulumi.Int(httpsPort),
                ToPort:     pulumi.Int(httpsPort),
                CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
            },
        },
    })
    if err != nil {
        return nil, err
    }

    // Create IAM roles for EKS
    eksRole, err := iam.NewRole(ctx, eksRoleName, &iam.RoleArgs{
        AssumeRolePolicy: pulumi.String(EksRoleAssumeRolePolicy),
    })
    if err != nil {
        return nil, err
    }

    // Attach policies to EKS role
    _, err = iam.NewRolePolicyAttachment(ctx, attachPolicy1Name, &iam.RolePolicyAttachmentArgs{
        Role:      eksRole.Name,
        PolicyArn: pulumi.String(eksServiceRolePolicy1),
    })
    if err != nil {
        return nil, err
    }

    _, err = iam.NewRolePolicyAttachment(ctx, attachPolicy2Name, &iam.RolePolicyAttachmentArgs{
        Role:      eksRole.Name,
        PolicyArn: pulumi.String(eksServiceRolePolicy2),
    })
    if err != nil {
        return nil, err
    }

    // Create EKS cluster
    eksCluster, err := eks.NewCluster(ctx, eksClusterResourceName, &eks.ClusterArgs{
        Name:    pulumi.String(eksClusterName),
        RoleArn: eksRole.Arn,
        VpcConfig: &eks.ClusterVpcConfigArgs{
            PublicAccessCidrs: pulumi.StringArray{
                pulumi.String("0.0.0.0/0"),
            },
            SecurityGroupIds: pulumi.StringArray{
                clusterSg.ID().ToStringOutput(),
            },
            SubnetIds: pulumi.StringArray{
                publicSubnet.ID(),
                privateSubnet.ID(),
            },
        },
    })
    if err != nil {
        return nil, err
    }

    // Create IAM role for Node Group
    nodeGroupRole, err := iam.NewRole(ctx, nodeGroupRoleName, &iam.RoleArgs{
        AssumeRolePolicy: pulumi.String(NodeGroupRoleAssumeRolePolicy),
    })
    if err != nil {
        return nil, err
    }

    // Attach policies to Node Group role
    _, err = iam.NewRolePolicyAttachment(ctx, nodeGroupPolicyAttachment1, &iam.RolePolicyAttachmentArgs{
        Role:      nodeGroupRole.Name,
        PolicyArn: pulumi.String(nodeGroupPolicy1),
    })
    if err != nil {
        return nil, err
    }

    _, err = iam.NewRolePolicyAttachment(ctx, nodeGroupPolicyAttachment2, &iam.RolePolicyAttachmentArgs{
        Role:      nodeGroupRole.Name,
        PolicyArn: pulumi.String(nodeGroupPolicy2),
    })
    if err != nil {
        return nil, err
    }

    _, err = iam.NewRolePolicyAttachment(ctx, nodeGroupPolicyAttachment3, &iam.RolePolicyAttachmentArgs{
        Role:      nodeGroupRole.Name,
        PolicyArn: pulumi.String(nodeGroupPolicy3),
    })
    if err != nil {
        return nil, err
    }

    // Create Node Group
    _, err = eks.NewNodeGroup(ctx, nodeGroupName, &eks.NodeGroupArgs{
        ClusterName:   eksCluster.Name,
        NodeRoleArn:   nodeGroupRole.Arn,
        SubnetIds:     pulumi.StringArray{publicSubnet.ID(), privateSubnet.ID()},
        InstanceTypes: pulumi.StringArray{pulumi.String(instanceType)},
        ScalingConfig: &eks.NodeGroupScalingConfigArgs{
            DesiredSize: pulumi.Int(desiredSize),
            MinSize:     pulumi.Int(minSize),
            MaxSize:     pulumi.Int(maxSize),
        },
    })
    if err != nil {
        return nil, err
    }

    // Create ACM Certificate
    certificate, err := acm.NewCertificate(ctx, "cert", &acm.CertificateArgs{
        DomainName:       pulumi.String(domainAddress),
        ValidationMethod: pulumi.String("DNS"),
    })
    if err != nil {
        return nil, err
    }

	return &Infrastructure{
		vpc:    vpc,
		server: eksCluster,
		group:  clusterSg,
		cert:   certificate,
	}, nil
}

// Infrastructure represents resources created for the EKS cluster.
type Infrastructure struct {
	vpc    *ec2.Vpc
	server *eks.Cluster
	group  *ec2.SecurityGroup
	cert   *acm.Certificate
}

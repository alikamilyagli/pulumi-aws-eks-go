package main

import (
	"sync"
	"testing"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mocks int

// Create mocks for Pulumi resources.
func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	outputs := args.Inputs.Mappable()
	switch args.TypeToken {
        case "aws:ec2/vpc:Vpc":
            outputs["id"] = "vpc-12345"
        case "aws:ec2/subnet:Subnet":
            outputs["id"] = "subnet-12345"
        case "aws:ec2/securityGroup:SecurityGroup":
            outputs["id"] = "sg-12345"
        case "aws:iam/role:Role":
            outputs["arn"] = "arn:aws:iam::123456789012:role/eks-role"
        case "aws:eks/cluster:Cluster":
            outputs["name"] = eksClusterName
            outputs["endpoint"] = "https://example.com"
            outputs["arn"] = "arn:aws:eks:us-west-1:123456789012:cluster/my-eks-cluster"
        case "aws:acm/certificate:Certificate":
            outputs["arn"] = "arn:aws:acm:us-west-1:123456789012:certificate/cert-12345"
	}
	return args.Name + "_id", resource.NewPropertyMapFromMap(outputs), nil
}

// Call mocks for Pulumi calls.
func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	outputs := map[string]interface{}{}
	if args.Token == "aws:ec2/getAmi:getAmi" {
		outputs["architecture"] = "x86_64"
		outputs["id"] = "ami-0eb1f3cdeeb8eed2a"
	}
	return resource.NewPropertyMapFromMap(outputs), nil
}

// TestInfrastructure runs unit tests to verify AWS infrastructure setup.
func TestInfrastructure(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		// Use SetupInfrastructure function from main.go
		infra, err := SetupInfrastructure(ctx)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		wg.Add(3)

		// Test if the EKS cluster has the correct name
		pulumi.All(infra.server.URN(), infra.server.Name).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			name := all[1].(string)

			assert.Equal(t, eksClusterName, name, "Cluster name does not match on server %v", urn)
			wg.Done()
			return nil
		})

		// Test if SSH port 22 is not open to the internet
		pulumi.All(infra.group.URN(), infra.group.Ingress).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			ingress := all[1].([]ec2.SecurityGroupIngress)

			for _, i := range ingress {
				openToInternet := false
				for _, b := range i.CidrBlocks {
					if b == "0.0.0.0/0" {
						openToInternet = true
						break
					}
				}

				assert.Falsef(t, i.FromPort == sshPort && openToInternet, "illegal SSH port 22 open to the Internet (CIDR 0.0.0.0/0) on group %v", urn)
			}

			wg.Done()
			return nil
		})

		// Test if ACM certificate has the correct domain name
		pulumi.All(infra.cert.URN(), infra.cert.DomainName).ApplyT(func(all []interface{}) error {
			urn := all[0].(pulumi.URN)
			domainName := all[1].(string)

			assert.Equal(t, domainAddress, domainName, "Domain name does not match on certificate %v", urn)
			wg.Done()
			return nil
		})

		wg.Wait()
		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))
	assert.NoError(t, err)
}

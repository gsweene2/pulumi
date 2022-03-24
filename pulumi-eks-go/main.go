package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		prefix := "pulumi-eks-go"

		// Resource: VPC
		// Purpose: Amazon Virtual Private Cloud (Amazon VPC) enables you to launch AWS resources into a virtual network that you've defined.
		// Docs: https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html

		// VPC CIDR
		cidrBlock := "10.0.0.0/16"

		// VPC Args
		vpcArgs := &ec2.VpcArgs{
			CidrBlock:          pulumi.String(cidrBlock),
			EnableDnsHostnames: pulumi.Bool(true),
			InstanceTenancy:    pulumi.String("default"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-vpc"),
			},
		}

		// VPC
		vpc, err := ec2.NewVpc(ctx, prefix+"-vpc", vpcArgs)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		// Resource: Subnets
		// Purpose: A subnet is a range of IP addresses in your VPC.
		// Docs: https://docs.aws.amazon.com/vpc/latest/userguide/configure-subnets.html

		// 3 Private Subnets
		privSubnet1, err := ec2.NewSubnet(ctx, prefix+"-priv-subnet-1", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.1.0/24"),
			AvailabilityZone: pulumi.String("us-east-2a"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-priv-subnet-1"),
			},
		})
		if err != nil {
			return err
		}
		privSubnet2, err := ec2.NewSubnet(ctx, prefix+"-priv-subnet-2", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.2.0/24"),
			AvailabilityZone: pulumi.String("us-east-2b"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-priv-subnet-2"),
			},
		})
		if err != nil {
			return err
		}
		privSubnet3, err := ec2.NewSubnet(ctx, prefix+"-priv-subnet-3", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.3.0/24"),
			AvailabilityZone: pulumi.String("us-east-2c"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-priv-subnet-3"),
			},
		})
		if err != nil {
			return err
		}

		// 3 Public Subnets
		pubSubnet1, err := ec2.NewSubnet(ctx, prefix+"-pub-subnet-1", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.4.0/24"),
			AvailabilityZone: pulumi.String("us-east-2a"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-pub-subnet-1"),
			},
		})
		if err != nil {
			return err
		}
		pubSubnet2, err := ec2.NewSubnet(ctx, prefix+"-pub-subnet-2", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.5.0/24"),
			AvailabilityZone: pulumi.String("us-east-2b"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-pub-subnet-2"),
			},
		})
		if err != nil {
			return err
		}
		pubSubnet3, err := ec2.NewSubnet(ctx, prefix+"-pub-subnet-3", &ec2.SubnetArgs{
			VpcId:            vpc.ID(),
			CidrBlock:        pulumi.String("10.0.6.0/24"),
			AvailabilityZone: pulumi.String("us-east-2c"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-pub-subnet-3"),
			},
		})
		if err != nil {
			return err
		}

		// Resource: Elastic IP
		// Purpose: An Elastic IP address is a static IPv4 address designed for dynamic cloud computing.
		// Docs: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html

		// EIP for NAT GW
		eip1, err := ec2.NewEip(ctx, prefix+"-eip1", &ec2.EipArgs{
			Vpc: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Resource: NAT Gateway
		// Purpose: A NAT gateway is a Network Address Translation (NAT) service.
		// Docs: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html

		// NAT Gateway with EIP
		natGw1, err := ec2.NewNatGateway(ctx, prefix+"-nat-gw-1", &ec2.NatGatewayArgs{
			AllocationId: eip1.ID(),
			// NAT must reside in public subnet for private instance internet access
			SubnetId: pubSubnet1.ID(),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-nat-gw-1"),
			},
		})
		if err != nil {
			return err
		}

		// Resource: Internet Gateway
		// Purpose: An internet gateway is a horizontally scaled, redundant, and highly available VPC component that allows communication between your VPC and the internet.
		// Docs: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Internet_Gateway.html

		// IGW for the Public Subnets
		igw1, err := ec2.NewInternetGateway(ctx, prefix+"-gw", &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-gw"),
			},
		})
		if err != nil {
			return err
		}

		// Resource: Route Tables
		// Purpose: A route table contains a set of rules, called routes, that determine where network traffic from your subnet or gateway is directed.
		// Docs: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html

		// Private Route Table for Private Subnets
		privateRouteTable, err := ec2.NewRouteTable(ctx, prefix+"-rtb-private-1", &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					// To Internet via NAT
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: natGw1.ID(),
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-rtb-private-1"),
			},
		})
		if err != nil {
			return err
		}

		// Public Route Table for Public Subnets
		publicRouteTable, err := ec2.NewRouteTable(ctx, prefix+"-rtb-public-1", &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				// To Internet via IGW
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: igw1.ID(),
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-rtb-public-1"),
			},
		})
		if err != nil {
			return err
		}

		// Associate Private Subs with Private Route Tables
		_, err = ec2.NewRouteTableAssociation(ctx, prefix+"-rtb-assoc-priv-1", &ec2.RouteTableAssociationArgs{
			SubnetId:     privSubnet1.ID(),
			RouteTableId: privateRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, prefix+"-rtb-assoc-priv-2", &ec2.RouteTableAssociationArgs{
			SubnetId:     privSubnet2.ID(),
			RouteTableId: privateRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, prefix+"-rtb-assoc-priv-3", &ec2.RouteTableAssociationArgs{
			SubnetId:     privSubnet3.ID(),
			RouteTableId: privateRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		// Associate Public Subs with Public Route Tables
		_, err = ec2.NewRouteTableAssociation(ctx, prefix+"-rtb-assoc-pub-1", &ec2.RouteTableAssociationArgs{
			SubnetId:     pubSubnet1.ID(),
			RouteTableId: publicRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, prefix+"-rtb-assoc-pub-2", &ec2.RouteTableAssociationArgs{
			SubnetId:     pubSubnet2.ID(),
			RouteTableId: publicRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRouteTableAssociation(ctx, prefix+"-rtb-assoc-pub-3", &ec2.RouteTableAssociationArgs{
			SubnetId:     pubSubnet3.ID(),
			RouteTableId: publicRouteTable.ID(),
		})
		if err != nil {
			return err
		}

		// Resource: IAM Role
		// Purpose: An IAM role is an IAM identity that you can create in your account that has specific permissions.
		// Docs: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html

		// IAM Role for EKS
		eksRole, err := iam.NewRole(ctx, prefix+"eks-iam-eksRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
		    "Version": "2008-10-17",
		    "Statement": [{
		        "Sid": "",
		        "Effect": "Allow",
		        "Principal": {
		            "Service": "eks.amazonaws.com"
		        },
		        "Action": "sts:AssumeRole"
		    }]
		}`),
		})
		if err != nil {
			return err
		}
		eksPolicies := []string{
			"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
			"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
		}
		for i, eksPolicy := range eksPolicies {
			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf(prefix+"rpa-%d", i), &iam.RolePolicyAttachmentArgs{
				PolicyArn: pulumi.String(eksPolicy),
				Role:      eksRole.Name,
			})
			if err != nil {
				return err
			}
		}

		// Resource: Security Groups
		// Purpose: A security group acts as a virtual firewall, controlling the traffic that is allowed to reach and leave the resources that it is associated with.
		// Docs: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html

		// Create a Security Group that we can use to actually connect to our cluster
		clusterSg, err := ec2.NewSecurityGroup(ctx, prefix+"-cluster-sg", &ec2.SecurityGroupArgs{
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
					FromPort:   pulumi.Int(80),
					ToPort:     pulumi.Int(80),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-cluster-sg"),
			},
		})
		if err != nil {
			return err
		}

		// Resource: EKS Cluster
		// Purpose: The Amazon EKS control plane & Amazon EKS nodes that are registered with the control plane
		// Docs: https://docs.aws.amazon.com/eks/latest/userguide/clusters.html

		// EKS cluster
		var subnetIds pulumi.StringArray

		cluster, err := eks.NewCluster(ctx, prefix+"-cluster", &eks.ClusterArgs{
			RoleArn: pulumi.StringOutput(eksRole.Arn),
			VpcConfig: &eks.ClusterVpcConfigArgs{
				PublicAccessCidrs: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				SecurityGroupIds: pulumi.StringArray{
					clusterSg.ID().ToStringOutput(),
				},
				SubnetIds: pulumi.StringArray(
					append(
						subnetIds,
						privSubnet1.ID(),
						privSubnet2.ID(),
						privSubnet3.ID(),
						pubSubnet1.ID(),
						pubSubnet2.ID(),
						pubSubnet3.ID(),
					)),
			},
		})
		if err != nil {
			return err
		}

		// Resource: IAM Role
		// Purpose: An IAM role is an IAM identity that you can create in your account that has specific permissions.
		// Docs: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html

		// Create the EC2 NodeGroup Role
		nodeGroupRole, err := iam.NewRole(ctx, prefix+"nodegroup-iam-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
		    "Version": "2012-10-17",
		    "Statement": [{
		        "Sid": "",
		        "Effect": "Allow",
		        "Principal": {
		            "Service": "ec2.amazonaws.com"
		        },
		        "Action": "sts:AssumeRole"
		    }]
		}`),
		})
		if err != nil {
			return err
		}
		nodeGroupPolicies := []string{
			"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
			"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
			"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
		}
		for i, nodeGroupPolicy := range nodeGroupPolicies {
			_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf(prefix+"ngpa-%d", i), &iam.RolePolicyAttachmentArgs{
				Role:      nodeGroupRole.Name,
				PolicyArn: pulumi.String(nodeGroupPolicy),
			})
			if err != nil {
				return err
			}
		}

		// Resource: EKS Node Groups
		// Purpose: Amazon EKS managed node groups automate the provisioning and lifecycle management of nodes (Amazon EC2 instances) for Amazon EKS Kubernetes clusters.
		// Docs: https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html

		var privsubnetIds pulumi.StringArray

		// Node Group 1
		_, err = eks.NewNodeGroup(ctx, prefix+"-worker-group-1", &eks.NodeGroupArgs{
			ClusterName: cluster.Name,
			SubnetIds: pulumi.StringArray(
				append(
					privsubnetIds,
					privSubnet1.ID(),
					privSubnet2.ID(),
					privSubnet3.ID(),
				)),
			NodeRoleArn: pulumi.StringInput(nodeGroupRole.Arn),
			ScalingConfig: &eks.NodeGroupScalingConfigArgs{
				DesiredSize: pulumi.Int(1),
				MinSize:     pulumi.Int(1),
				MaxSize:     pulumi.Int(1),
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-worker-group-1"),
			},
		})
		if err != nil {
			return err
		}

		// Node Group 2
		_, err = eks.NewNodeGroup(ctx, prefix+"-worker-group-2", &eks.NodeGroupArgs{
			ClusterName: cluster.Name,
			SubnetIds: pulumi.StringArray(
				append(
					privsubnetIds,
					privSubnet1.ID(),
					privSubnet2.ID(),
					privSubnet3.ID(),
				)),
			NodeRoleArn: pulumi.StringInput(nodeGroupRole.Arn),
			ScalingConfig: &eks.NodeGroupScalingConfigArgs{
				DesiredSize: pulumi.Int(1),
				MinSize:     pulumi.Int(1),
				MaxSize:     pulumi.Int(1),
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String(prefix + "-worker-group-2"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}

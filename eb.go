package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	eb "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/elb"
)

type EBApp struct {
	App  *eb.ApplicationDescription
	Envs []EBEnv
}

type EBEnv struct {
	Env       *eb.EnvironmentDescription
	EnvRes    *eb.EnvironmentResourceDescription
	Instances []*ec2.Instance
	ELBs      []*elb.LoadBalancerDescription
}

func listAllEB() ([]*ec2.Instance, error) {

	ebclient := eb.New(sess)
	ec2client := ec2.New(sess)
	elbclient := elb.New(sess)

	apps, err := ebclient.DescribeApplications(&eb.DescribeApplicationsInput{})
	if err != nil {
		return nil, err
	}

	for _, app := range apps.Applications {

		envs, err := ebclient.ComposeEnvironments(&eb.ComposeEnvironmentsInput{ApplicationName: app.ApplicationName})

		if err != nil {
			return nil, err
		}
		for _, env := range envs.Environments {
			envResorces, err := ebclient.DescribeEnvironmentResources(&eb.DescribeEnvironmentResourcesInput{EnvironmentId: env.EnvironmentId})
			if err != nil {
				return nil, err
			}

			evEnv := EBEnv{
				Env:    env,
				EnvRes: envResorces.EnvironmentResources,
			}

			{ // instance
				var instanceIDs []*string
				for _, instance := range envResorces.EnvironmentResources.Instances {
					instanceIDs = append(instanceIDs, instance.Id)
				}

				instanceSearchCond := &ec2.DescribeInstancesInput{
					InstanceIds: instanceIDs,
					MaxResults:  aws.Int64(1000),
				}

				instances, err := ec2client.DescribeInstances(instanceSearchCond)

				for _, resv := range instances.Reservations {
					for _, instance := range resv.Instances {
						evEnv.Instances = append(evEnv.Instances, instance)
					}
				}

				for instances.NextToken != nil {
					instanceSearchCond.NextToken = instances.NextToken
					instances, err = ec2client.DescribeInstances(instanceSearchCond)
					if err != nil {
						return nil, err
					}

					for _, resv := range instances.Reservations {
						for _, instance := range resv.Instances {
							evEnv.Instances = append(evEnv.Instances, instance)
						}
					}
				}
			}

			// elb
			{
				var lbnames []*string
				for _, elbInfo := range envResorces.EnvironmentResources.LoadBalancers {
					lbnames = append(lbnames, elbInfo.Name)
				}

				//** "page size" unsupported now **
				elbOut, err := elbclient.DescribeLoadBalancers(&elb.DescribeLoadBalancersInput{LoadBalancerNames: lbnames})
				if err != nil {
					return nil, err
				}
				evEnv.ELBs = elbOut.LoadBalancerDescriptions
			}

		}
	}
	return nil, nil

}

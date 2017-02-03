package awstool

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	eb "github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/go-xweb/log"
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

func ListEB(region string, filterAppNames []*string) ([]EBApp, error) {

	sess := session.New(&aws.Config{Region: aws.String(region)})

	ebclient := eb.New(sess)
	ec2client := ec2.New(sess)
	elbclient := elb.New(sess)

	apps, err := ebclient.DescribeApplications(&eb.DescribeApplicationsInput{
		ApplicationNames: filterAppNames,
	})

	if err != nil {
		log.Error(" ad DescribeApplications %s", err.Error())
		return nil, err
	}

	var ebApps []EBApp
	for _, app := range apps.Applications {
		ebApp := EBApp{App: app}
		envs, err := ebclient.DescribeEnvironments(&eb.DescribeEnvironmentsInput{ApplicationName: app.ApplicationName})
		if err != nil {
			log.Error(" error ComposeEnvironments %s", err.Error())
			return nil, err
		}

		var ebEnvs []EBEnv
		for _, env := range envs.Environments {
			envResorces, err := ebclient.DescribeEnvironmentResources(&eb.DescribeEnvironmentResourcesInput{EnvironmentId: env.EnvironmentId})
			if err != nil {
				log.Error(" error DescribeEnvironmentResources%s", err.Error())
				return nil, err
			}

			ebEnv := EBEnv{
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
						ebEnv.Instances = append(ebEnv.Instances, instance)
					}
				}

				for instances.NextToken != nil {
					instanceSearchCond.NextToken = instances.NextToken
					instances, err = ec2client.DescribeInstances(instanceSearchCond)
					if err != nil {
						log.Error(" error DescribeInstances %s", err.Error())
						return nil, err
					}

					for _, resv := range instances.Reservations {
						for _, instance := range resv.Instances {
							ebEnv.Instances = append(ebEnv.Instances, instance)
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
					log.Error(" error DescribeLoadBalancers %s", err.Error())
					return nil, err
				}
				ebEnv.ELBs = elbOut.LoadBalancerDescriptions
			}

			ebEnvs = append(ebEnvs, ebEnv)
		}

		ebApp.Envs = ebEnvs

		ebApps = append(ebApps, ebApp)
	}
	return ebApps, nil

}

## Install
```
go get github.com/tacogips/awstool
```

## Usage
- put config file at $HOME/.awstool (sample file is at cmd/awstool/awstool_config.yaml)


## elastic beanstalk
### list

```
awstool eb list
```
### output obj struct(for template)

```
type EBApp struct {
	App  *eb.ApplicationDescription
	Envs []EBEnv
}

type EBEnv struct {
	Env                 *eb.EnvironmentDescription
	EnvRes              *eb.EnvironmentResourceDescription
	InstanceAndStatuses []InstanceAndStatus
	ELBs                []*elb.LoadBalancerDescription
}

type InstanceAndStatus struct {
	Instance *ec2.Instance
	Status   *ec2.InstanceStatus
}
```


## s3
### list

```
awstool s3 list --pre="your/key"
```
### output obj struct(for template)

```
type S3Files struct {
	Bucket  string
	SumSize int64
	Count   int
	Objects []*s3.Object
}
```

### dl

```
awstool s3 dl --pre="your/key"
```


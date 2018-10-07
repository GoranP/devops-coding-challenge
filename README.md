# Engagetech DevOps Coding Test

# Solution
Idea for solution is to create infrastructure is base on initial drawing (https://github.com/GoranP/devops-coding-challenge/blob/master/src/scheme.jpg)

* Infrastructure is build with cli tools and wraped up in bash scripts.
* Deployment script is written in bash.
* Simple web service that returns current time is written in golang, and packed in docker.
* Health check script is written in bash.

On AWS following will be built:
 * Elastic Container Service (ECS) cluster (FARGATE)
 * ECS service and tasks
 * Application Load Balancer connected with ECS
 * VPC
 * CloudWatch logs of application
 * all necessary security groups and roles for ECS cluster to work

With this infrastructure service is highly available and secured behind AWS secure groups.

# Installation of prerequisites
Solution for challenge requires aws cli tools and latest docker installed on the system. 
It is tested on Ubuntu 18.04. Best would be to spawn new aws ubuntu 18.04  instance and run this solution there.

Install and configure latest AWS cli tools:
```sh
$ sudo snap install --channel edge aws-cli --classic
$ aws configure
AWS Access Key ID [****************]: <enter valid AWS acces key>
AWS Secret Access Key [****************]: <enter valid AWS secret key>
Default region name [eu-west-1]: eu-west-1
Default output format [json]: json
```
## Docker
Install latest docker following official procedure:
https://docs.docker.com/install/linux/docker-ce/ubuntu/

Take a special attention for enabling usage of docker without sudo:
https://docs.docker.com/install/linux/linux-postinstall/


Install helper AWS tool for Elastic Cloud Services:

```sh
$ sudo curl -o /usr/local/bin/ecs-cli https://s3.amazonaws.com/amazon-ecs-cli/ecs-cli-linux-amd64-latest
$ sudo chmod +x /usr/local/bin/ecs-cli
```
ecs-cli tool should use previously configured aws profile and it's credentilas. For more information on non-default profiles of ecs  and aws cli please read: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ECS_CLI_Configuration.html

**All tools locations on the system should be in shell path!**
# Assumptions
In this solution there is one assumtion: docker repository is already created and exists on AWS on the following uri: 434227579614.dkr.ecr.eu-west-1.amazonaws.com/gimmetime

Although it's trivial to create repo, I decided to leave existing one to speed up docker image upload during development, and repo is part of infrastructure that is usually not that easily deleted and recreated. There is a code and images for production. :)

Solution should be configured with API credentials of my aws account that will be supplied over secure channel.
Access to AWS web console will also be provided over secure channel.

# Usage
Clone this fork
```sh
$ git clone https://github.com/GoranP/devops-coding-challenge.git
```

In shell change directory to src:
```sh
$ cd /path/to/fork/src
```
To get repo credentials run:
```sh
$ $(aws ecr get-login --no-include-email --region eu-west-1) 
```

Verify that credentials are configured correctly:
```sh
$ ecs-cli images
REPOSITORY NAME     TAG                 IMAGE DIGEST                                                              PUSHED AT           SIZE                
gimmetime           latest              sha256:4a1efa87fdbb7ce22166e587b2e042cd53b6b26f50c77f090e3b687635f882e7   4 hours ago         300 MB              
gimmetime           <none>              sha256:aca61c94363a92c56eb8a177dce98155b08463bb549f55575db75a3dc9c120a3   2 days ago          300 MB         
$ 
```

If you get list of images and repos your credentials are configured good.

Build and deploy new docker with latest code:
```sh
$ cd gimmetime
$ ./build-and-push
```

Go back to _src_ directory:
```sh
cd ..
```

Run script _build-infrastructure_
```sh
$ ./build-infrastructure
```

After 3-5 minutes infrastructure on AWS should be created.
Follow instructions at the end.

Example of success info:
```
$ ./build-infrastructure
Configuring ecs-cli cluster file
INFO[0000] Saved ECS CLI cluster configuration default. 
Creating ecsTaskExecutionRole and ecsServiceExecutionRole roles
Creating ecs cluster...
INFO[0000] Created cluster                               cluster=engage-infra region=eu-west-1
INFO[0001] Waiting for your cluster resources to be created... 
INFO[0001] Cloudformation stack status                   stackStatus=CREATE_IN_PROGRESS
Creating security groups for ECS and ALB
Creating Application Load Balancer...
Creating and starting ECS service and task with ALB connected
INFO[0000] Using ECS task definition                     TaskDefinition="engage-infra:27"
INFO[0001] Created Log Group engage-infra-grp in eu-west-1 
INFO[0001] Created an ECS service                        service=engage-infra taskDefinition="engage-infra:27"
WARN[0001] Failed to create log group engage-infra-grp in eu-west-1: The specified log group already exists 
INFO[0001] Updated ECS service successfully              desiredCount=1 serviceName=engage-infra
INFO[0017] (service engage-infra) has started 1 tasks: (task 87b13940-0e8d-40e2-a2e9-a0adfa182242).  timestamp="2018-10-07 13:33:19 +0000 UTC"
INFO[0093] Service status                                desiredCount=1 runningCount=1 serviceName=engage-infra
INFO[0093] ECS Service has reached a stable state        desiredCount=1 runningCount=1 serviceName=engage-infra

ECS cluster with load balancer successfully created!

To verify runnig container run: curl http://infra-ecs-alb-1119285097.eu-west-1.elb.amazonaws.com/now
Service returns UNIX format of time, number of seconds since 1.1.1970
To get meaningful datetime run following:
date --date @$(curl --silent http://infra-ecs-alb-1119285097.eu-west-1.elb.amazonaws.com/now)

To scale cluster up and down run following (with desired number of tasks in cluster):
aws ecs update-service --cluster engage-infra --service engage-infra --desired-count 4
Alternativley login to AWS console on web and update cluster service and set desired count


*** Cleanup helper
To delete all resources ensure that cluster is scaled to 0 (first commmand) and then all other commands:
--------------------------------------------------------------------------------------------------------

  aws ecs update-service --cluster engage-infra --service engage-infra --desired-count 0
  aws ecs delete-service --cluster engage-infra --service engage-infra

  aws elbv2 delete-load-balancer --load-balancer-arn arn:aws:elasticloadbalancing:eu-west-1:434227579614:loadbalancer/app/infra-ecs-alb/723a593a9eb1feb0
  aws elbv2 delete-target-group --target-group-arn arn:aws:elasticloadbalancing:eu-west-1:434227579614:targetgroup/infra-gimmetime/3afa2e3988b86bae

  aws ec2 delete-security-group --group-id sg-02e71a0bb9e88a251
  aws ec2 delete-security-group --group-id sg-0a9d54c51e68c3b9d

  ecs-cli down --force --cluster engage-infra	

  aws iam  detach-role-policy --role-name ecsTaskExecutionRole --policy-arn arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
  aws iam  delete-role --role-name ecsTaskExecutionRole
  aws iam  detach-role-policy --role-name ecsServiceExecutionRole --policy-arn arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole
  aws iam  delete-role --role-name ecsServiceExecutionRole

  aws logs delete-log-group --log-group-name engage-infra-grp

```

To check cluster and service availability run curl command that is provided in success info.
To scale cluster from one task to 4 run:
```sh
$ aws ecs update-service --cluster engage-infra --service engage-infra --desired-count 4
```
To verify time difference, as requested by challenge, run health check script:
```sh
$ ./health-check http://<dns_of_alb_DNS_visible_in_info_after_infra_build>/now
0
```
To get DNS of ALB you can also log in to AWS console, locate ALB in Ireland region and copy paste DNS name. 

To view logs of service in CloudWatch visit AWS console: https://eu-west-1.console.aws.amazon.com/cloudwatch/home?region=eu-west-1#logs:


To deploy latest image first change someting in main.go and build new docker image in gimmetime directory:
```sh
$ cd gimmetime
$ nano main.go 
$ ./build-and-push
```

And the run _deploy-latest_ script in src direcotry:
```sh
$ cd ..
$ ./deploy-latest
```

Deploy will create new revision of task definition on ECS, and will update ecs service with new task revision.
This will instruct ECS to automatically trigger deployment of latest docker image, and automatically drains old images before stopping them. Essentially this automatical rolling deployment.


# Further improvements
Although solution is higly available behind ALB, and secure behind AWS secure groups it would be good to secure it further.
* encrypt traffic, use https on with SSL certificate ALB
* configure automatic scaling of ECS cluster. Currently it is manual.
* use CloudTrail to audit all actions made on AWS infrastructure by aws accounts
* implement AWS WAF
* use Amazon Guard Duty for threat intelligence


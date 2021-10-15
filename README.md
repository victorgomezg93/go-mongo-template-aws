# Terraform-aws-go

---
DEPLOYMENT
----

This project is build to provide an infraestructure in AWS to deploy a Go server.

In case you want to deploy it you have to do the following:

First in your command line you have to configure the access in aws

```sh
aws configure
```

Then you have an initialization script, it's called "init.bat", this script prompts you and asks for the registryID, when you give this ID the script builds and tag the docker image, create an ECR repository and pushed the imatge into the ECR, also the script checks for terraform configuration files because it's an init script so in case we have an tfstate or lock the script remove this files

```sh
init.bat
```

Once initialized we have to change the registryID in variables.tf with your own ID, here:

```sh
variable "app_image" {
  description = "Docker image to run in the ECS cluster"
  default     = "880231462042.dkr.ecr.us-east-1.amazonaws.com/go-ecs-app-repo:latest"
}
```

 and also create the certificate in AWS manually using the certificates located in /app, in the certificate manager you have to use public.crt, private.key and rootCA.crt in this order to create the certificate and copy the arn and paste it in ALB here:
 
```sh
resource "aws_alb_listener" "front_end" {
  load_balancer_arn = aws_alb.main.id
  port              = var.app_port
  protocol          = "HTTPS"
  certificate_arn    = "arn:aws:acm:us-east-1:880231462042:certificate/8730eba6-f34c-46d0-921b-0a460ee1a181"
```
Then we can go to the infraestructure folder and use:

```sh
terraform init
terraform apply --auto-approve
```

The result should be an alb hostname accesible by https with a healtcheck in /health to check one database.

I'm using a sandbox with automated shutdown so this urls won't be accesible but in my case after deployment I recieve a message with this alb hostname: alb_hostname = "myapp-load-balancer-920324680.us-east-1.elb.amazonaws.com"

If i go to this urls I would be able to test the certificate and also information about one database not deployed
 
https://myapp-load-balancer-920324680.us-east-1.elb.amazonaws.com
https://myapp-load-balancer-920324680.us-east-1.elb.amazonaws.com/health

INFRAESTRUCTURE EXPLANATION
----

For our deployment in the cloud we are using terraform.

We have a Cluster with a service deployed with 3 different tasks, this tasks are our go image and they are communicated with a load balancer in port 443, also this load balancer is using a certificate deployed before.
Using the application load balancer we can distribute  network traffic and information flows across multiple servers, a load balancer ensures no single server bears too much demand. This improves application responsiveness and availability, enhances user experiences, and can protect from distributed denial-of-service (DDoS) attacks.

Also the service is autoscaled depending of the CPU utilization, minimum we run 3 tasks but we can scale up to 6 if the cpu utilizations goes >= 85, if the trafic is reduced and we find that we have  <= 10 cpu utlization we remove one task to avoid cost.

We are using ECS because is really powerful compared with EC2 and we have an internet gateway for the public subnets with public and private subnets where the containers is located. We are deploying it on us-east-1 because my sandbox is located in this region.

Our container is a go server is running with https protocol in port 443, the aws credentials are gathered from the credentials aws file in $HOME/.aws/credentials.

Map explanation (hard mode)

![alt text](https://github.com/victorgomezg93/terraform-aws-go/blob/main/graph.png?raw=true)

Better view

![alt text](https://github.com/victorgomezg93/terraform-aws-go/blob/main/diagram.PNG?raw=true)

CODE
----
We have a docker to deploy the application and to implement some api we have in the /health path a healtcheck checking if our sqlite3 database is connected or not, I didn't had time to implement the database so always return the same response in json:

```sh
{"status":"down","details":{"database":{"status":"down","timestamp":"2021-10-05T23:24:08.425492059Z","error":"Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub"},"search":{"status":"down","timestamp":"2021-10-05T23:24:06.309964687Z","error":"this makes the check fail"}}}
```

The only changing thing is the timestamp.

GUIDE TO CREATE PRIVATE/PUBLIC KEYS (not necessary, just as a reminder when the certificate expires)
----

First we created the CA

```sh
openssl genrsa -des3 -out rootCA.key 4096
```

If you want a non password protected key you can remove the -des3 option but I will use one password

Then I created the self sign root certificate

```sh
wsl openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -out rootCA.crt
```

I created the certificate key

```sh
openssl genrsa -out private.key 2048
```

And then the sign request to generate the certificate, it's important to provide names in this case to avoid problems with AWS

```sh
wsl openssl x509 -req -in mydomain.com.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out public.crt -days 500 -sha256
```

And then generated the certificate

```sh
wsl openssl x509 -req -in mydomain.com.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out public.crt -days 500 -sha256
```
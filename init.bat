set /P ACC_ID=Enter AWS RegistryID: 

aws ecr get-login-password | docker login --username AWS --password-stdin %ACC_ID%.dkr.ecr.us-east-1.amazonaws.com
aws ecr create-repository --repository-name go-ecs-app-repo
docker-compose build
docker-compose push
cd infraestructure
del .terraform.lock.hcl
del terraform.tfstate
del terraform.tfstate.backup
terraform init
terraform apply --auto-approve
cd ..
aws ecs describe-task-definition --task-definition go-task  --query taskDefinition > task-definition.json
PAUSE
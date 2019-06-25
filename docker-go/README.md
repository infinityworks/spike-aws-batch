# Usage

## Build
Build the docker image and push to the AWS ECR repo
```bash
# log into AWS ECR for the eu-west region
$(aws ecr get-login --no-include-email --region eu-west-1)
# build the docker image spike-aws-batch
docker build --tag=spike-aws-batch .
# tag & push the image teh the AWS ECR repo
docker tag spike-aws-batch:latest <AWS_REPO_URI>:latest
docker push <AWS_REPO_URI>:latest
```

## Configuration

use a tool to assume an AWS IAM Role with MFA [GITHUB remind101/assume-role](https://github.com/remind101/assume-role) for the session

```bash
assume-role iw-sandpit
```

## Pre & Dry Run 

Ensure the target bucket has glacier files within it
```bash
aws s3api list-objects-v2 \
    --bucket $RESTORE_BUCKET \
    --query "Contents[?StorageClass=='GLACIER']" \
    --output json \
    | jq -r '.[].Key' > list.txt
```

set the ENV variable RESTORE_DRYRUN passed into docker to allow a dry run (no resore requests just output the commands it would have run)
```bash
export RESTORE_DRYRUN=true
```

## Execution
execute the docker run passing in the various environment variables
```bash
 docker run \
 -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
 -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
 -e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
 -e AWS_SECURITY_TOKEN=${AWS_SECURITY_TOKEN} \
 -e ASSUMED_ROLE=${ASSUMED_ROLE} \
glacier-restore-go -bucket=spike-aws-batch -verbose=true -dryRun=true
 ```
# Usage

## Build

```bash
docker build --tag=glacier-restore .
```

## Local Export

```bash
assume-role iw-sandpit
env | grep AWS_
```

## Run

```bash
 docker run \
 -e BUCKET=spike-aws-batch \
 -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
 -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
 -e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
 -e AWS_SECURITY_TOKEN=${AWS_SECURITY_TOKEN} \
 -e ASSUMED_ROLE=${ASSUMED_ROLE} \
 glacier-restore s3 ls
 ```
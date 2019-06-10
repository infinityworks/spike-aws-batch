# Usage

## Build

```bash
docker build --tag=glacier-restore .
```

## Local Export

```bash
assume-role iw-sandpit
env | grep AWS_
export RESTORE_BUCKET=spike-aws-batch RESTORE_THREADS=5 RESTORE_DAYS=1 RESTORE_TIER=Bulk
```

## Test

```bash
aws s3api list-objects-v2 \
    --bucket $BUCKET \
    --query "Contents[?StorageClass=='GLACIER']" \
    --output json \
    | jq -r '.[].Key' > list.txt
```

## Run

```bash
 docker run \
 -e RESTORE_BUCKET=spike-aws-batch \
 -e RESTORE_THREADS=5 \
 -e RESTORE_DAYS=1 \
 -e RESTORE_TIER=Bulk \
 -e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
 -e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
 -e AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
 -e AWS_SECURITY_TOKEN=${AWS_SECURITY_TOKEN} \
 -e ASSUMED_ROLE=${ASSUMED_ROLE} \
 glacier-restore s3 ls
 ```
#!/bin/bash 

# Fetch the list a glacier object in the bucket
aws s3api list-objects-v2 \
    --bucket $RESTORE_BUCKET \
    --query "Contents[?StorageClass=='GLACIER']" \
    --output json \
    | jq -r '.[].Key' > list.txt

# make the working dir and split the list of files into multiple sub-lists
mkdir $RESTORE_BUCKET
cd $RESTORE_BUCKET
split -l $RESTORE_THREADS ../list.txt

for f in x*
do
source /glacier/sub.sh $f &
done
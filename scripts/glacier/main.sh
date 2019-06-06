#!/bin/bash

# TODO 
# pass in the file split number
# pass in the restore Tier
# pass in the restore Days

# Fetch the list a glacier object in the bucket
aws s3api list-objects-v2 \
    --bucket $BUCKET \
    --query "Contents[?StorageClass=='GLACIER']" \
    --output text \
| awk '{print $2}' > list.txt

# make the working dir and split the list of files into multiple sub-lists
mkdir $BUCKET
cd /$BUCKET
split -l 1000 ../list.txt

for f in x*
do
~glacier/sub.sh $f &
done
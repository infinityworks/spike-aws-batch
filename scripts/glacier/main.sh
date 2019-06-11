#!/bin/bash 

function sub() {
  while read p; do
    if $RESTORE_DRYRUN
    then
      echo 'aws s3api restore-object --restore-request '\''{"Days":'$RESTORE_DAYS',"GlacierJobParameters":{"Tier":"'$RESTORE_TIER'"}}'\'' --bucket "'$RESTORE_BUCKET'" --key "'$p'"'
    else
      echo 'restoring: '$RESTORE_BUCKET'/'$p
      aws s3api restore-object --restore-request '{"Days":'$RESTORE_DAYS',"GlacierJobParameters":{"Tier":"'$RESTORE_TIER'"}}' --bucket "$RESTORE_BUCKET" --key "$p"
    fi
  done < /glacier/$RESTORE_BUCKET/$1
}

# Fetch the list a glacier object in the bucket
aws s3api list-objects-v2 \
    --bucket $RESTORE_BUCKET \
    --query "Contents[?StorageClass=='GLACIER']" \
    --output json \
    | jq -r '.[].Key' > list.txt

# make the working dir and split the list of files into multiple sub-lists
mkdir $RESTORE_BUCKET
cd $RESTORE_BUCKET
split -l $RESTORE_PERTHREAD ../list.txt

for f in x*
do
sub $f &
done

wait
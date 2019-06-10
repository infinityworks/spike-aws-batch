#!/bin/bash

while read p; do
  aws s3api restore-object --restore-request '{"Days":'$RESTORE_DAYS',"GlacierJobParameters":{"Tier":"'$RESTORE_TIER'"}}' --bucket "$RESTORE_BUCKET" --key "$p"
done < /glacier/$RESTORE_BUCKET/$1
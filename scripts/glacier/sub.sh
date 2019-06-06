#!/bin/bash
while read p; do
  echo 'aws s3api restore-object  --restore-request '{"Days":7,"GlacierJobParameters":{"Tier":"Bulk"}}' --bucket $BUCKET --key '"$p"
#   aws s3api restore-object /
#     --restore-request '{ /
#         "Days" : 7, /
#         "GlacierJobParameters" : {"Tier":"Bulk"} /
#     }' /
#     --bucket $BUCKET /
#     --key "$p" /
done < ~glacier/$BUCKET/$1
# spike-aws-batch

An application, runtime and infrastructure to investigate AWS Batch
This project enables an AWS Batch to be stoodup that can be executed to make glacier restore requests for all files in a given S3 bucket that are in the Glacier storage class

## application

A Golang application that accepts:
- -bucket (string) the glacier bucket to restore
- [ -days ] (int) number of days to keep restoration alive for (default = 1)
- [ -tier ] (string) glacier restore tier, affects cost, Bulk/Standard/Expedited (default = Bulk)
- [ -region ] (strin) the default region to use (default = eu-west-1)
- [ -concurrency ] (int) number of concurrent aws s3 restore requests (default = 4)
- [ -dryRun ] passed to enable a dryrun
- [ -quiet ] passed to reduce logging output

## runtime

An Alpine Linux Docker with the compiled Golang application embedded as the Entrypoint

## infrastructure

A JSON CloudFormation template to stand up the AWS Batch resources and other to start the supporting/required AWS resources.
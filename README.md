# glacier-restore

A utility to make AWS S3 Glacier restore requests for all glaciered objects within an AWS S3 bucket

## Test

test from code with with verbose & dry run enabled.

```bash
go run main.go -bucket=spike-aws-batch -dryRun
```

## Build 

disabling cgo which gives us a static binary, compiling for linux.

```bash
CGO_ENABLED=0 GOOS=linux go build
```
 
## Run
execute with dry run enabled.


```bash
./glacier-restore -bucket spike-aws-batch -dryRun
```
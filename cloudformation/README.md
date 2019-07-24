# Example AWS CLI usage

```bash
aws batch submit-job \
--job-name test-batch-job-3 \
--job-queue HighPriority \
--job-definition aws-batch-example-task:7 \
--parameters bucket=spike-aws-batch,dryRun=true,verbose=true \
--region eu-west-1
```

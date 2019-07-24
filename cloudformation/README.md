# Example AWS CLI usage

```bash
aws batch submit-job \
--job-name test-batch-job-3 \
--job-queue HighPriority \
--job-definition aws-batch-example-task:14 \
--parameters bucket=spike-aws-batch \
--region eu-west-1
```

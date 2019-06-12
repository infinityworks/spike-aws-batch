provider "aws" {
  alias   = "global"
  version = "~> 2.12"
  region  = "eu-west-1"
  profile = "iw-sandpit"
}

terraform {
  backend "s3" {
    region  = "eu-west-1"
    key     = "projects/auth/state.json"
    bucket  = "spike-aws-batch.terraform.state"
    profile = "iw-sandpit"

    dynamodb_table = "spike-aws-batch.terraform.lock"
  }
}


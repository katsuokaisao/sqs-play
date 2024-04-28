terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "sql-play20240428014100918600000001"
    key    = "terraform/state"
    region = "ap-northeast-1"
  }
}

provider "aws" {}

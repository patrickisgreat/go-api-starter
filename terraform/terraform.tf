terraform {
  required_version = "~> 1.5"
  backend "s3" {
    key     = "services/go-api-starter"
    region  = "eu-central-1"
    encrypt = "true"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 3.75.0"
    }
    pagerduty = {
      source  = "pagerduty/pagerduty"
      version = "~> 3.30.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = ">= 4.10"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
  default_tags {
    tags = local.tags
  }
}


provider "newrelic" {}

provider "pagerduty" {}

provider "grafana" {}

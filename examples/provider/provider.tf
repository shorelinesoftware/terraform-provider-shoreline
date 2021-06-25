
terraform {
  required_providers {
    shoreline = {
      source  = "shoreline.io/terraform/shoreline"
      version = ">= 1.0.1"
    }
  }
}

provider "shoreline" {
  # provider configuration here
  #token = "xyz1.asdfj.asd3fas..."
  url = "https://test.us.api.shoreline-vm1.io"
  #url = "https://test.us.api.shoreline-test6.io"
  retries = 2
  debug   = true
}


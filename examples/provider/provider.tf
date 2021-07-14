
terraform {
  required_providers {
    shoreline = {
      source  = "shorelinesoftware/shoreline"
      version = ">= 1.0.6"
    }
  }
}

provider "shoreline" {
  # provider configuration here
  #token = "xyz1.asdfj.asd3fas..."
  url     = "https://acme.us.api.shoreline-cluster.io"
  retries = 2
  debug   = true
}


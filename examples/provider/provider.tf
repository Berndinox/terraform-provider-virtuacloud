terraform {
  required_providers {
    virtuacloud = {
      source  = "registry.opentofu.org/Berndinox/virtuacloud"
      version = "~> 0.1"
    }
  }
}

provider "virtuacloud" {
  api_key = var.virtuacloud_api_key
}
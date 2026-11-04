terraform {
  required_providers {
    tfcoremock = {
      source = "hashicorp/tfcoremock"
      version = "0.6.0"
    }
  }
}

provider "tfcoremock" {}

resource "tfcoremock_simple_resource" "simple_resource" {
  id = "192977d6-b169-4170-a9d4-ee1dcef7c6ea"
  integer = 1
}

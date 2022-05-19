#!/bin/bash

rm -rf ".terraform" ".terraform.lock.hcl" "terraform.tfstate" "terraform.tfstate.backup"
terraform init

terraform apply -auto-approve


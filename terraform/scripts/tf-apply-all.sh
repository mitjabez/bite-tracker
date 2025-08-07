#!/usr/bin/env bash
set -e

MY_DIR="$(dirname "$0")"

echo "Applying ECR first ..."
terraform apply -target aws_ecr_repository.bite_tracker

echo "Deploying demo app ..."
"$MY_DIR"/deploy-demo.sh
echo "Applying rest of infra ..."
terraform apply



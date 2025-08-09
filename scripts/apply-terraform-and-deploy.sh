#!/usr/bin/env bash
set -e

MY_DIR=$(dirname "$0")
MY_NAME=$(basename "$0")

usage() {
  cat <<EOF
Usage: ${MY_NAME}

Runs terraform apply and deploys the image. In general should be run just for infrastructure bootstrapping.
But you can also run it later to apply&deploy at the same time.


OPTIONS:
  -h --help   This help message.

EOF
}

main() {
  if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    usage
    exit 0
  fi

  echo "Applying ECR first ..."
  pushd "$MY_DIR/../terraform" > /dev/null
  terraform init
  terraform apply -target aws_ecr_repository.bite_tracker
  popd > /dev/null
  echo

  echo "Pushing image to ECR ..."
  "$MY_DIR"/deploy.sh push-image
  echo

  echo "Applying rest of infra ..."
  pushd "$MY_DIR/../terraform" > /dev/null
  terraform apply
  popd > /dev/null
}

main "$@"

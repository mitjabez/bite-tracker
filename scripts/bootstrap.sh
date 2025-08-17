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

die() {
  echo "ERROR: $*" >&2
  exit 1
}

requirements() {
  type -f terraform > /dev/null || die "terraform is not installed"
  type -f docker > /dev/null || die "docker is not installed"
}

main() {
  if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    usage
    exit 0
  fi

  requirements

  pushd "$MY_DIR/../terraform" > /dev/null
  [ -f terraform.tfvars ] || die "terraform.tfvars file is missing"

  echo "Bootstrapping the whole environment"
  echo "-----------------------------------"

  echo "Applying ECR and RDS first ..."
  terraform init
  terraform apply \
    -target aws_ecr_repository.bite_tracker
  popd > /dev/null


  echo "Pushing image to ECR ..."
  "$MY_DIR"/deploy.sh push-image
  echo

  echo "Applying rest of infra ..."
  pushd "$MY_DIR/../terraform" > /dev/null
	echo "Initial apply with DB role bootstrapping"
  terraform apply -var bootstrap_db_roles=true
  popd > /dev/null
}

main "$@"

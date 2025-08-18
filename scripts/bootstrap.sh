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

log() {
  echo "$@" >&2
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

  log "*** Bootsrapping Bite Tracker infrastructure ***"

	log
  log "** Applying ECR and RDS first **"
  terraform init
  terraform apply -target aws_ecr_repository.bite_tracker
  popd > /dev/null

	log
  log "* Pushing image to ECR *"
  "$MY_DIR"/deploy.sh push-image

	log
  log "** Applying rest of infra **"
  pushd "$MY_DIR/../terraform" > /dev/null
	log "* Initial apply with DB role bootstrapping permission *"
  terraform apply -var bootstrap_db_roles=true
	log "* Final apply to remove bootsrapping permissions *"
	terraform apply

	log
	log "*** All done! ***"

  popd > /dev/null
}

main "$@"

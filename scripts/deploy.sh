#!/usr/bin/env bash
set -e

MY_DIR="$(dirname "$0")"
MY_NAME=$(basename "$0")

usage() {
  cat <<EOF
Usage: ${MY_NAME} TARGET

Deploy or push image of bite-tracker to AWS App Runner


OPTIONS:
  TARGET      push-image: Only push the image
              trigger-deploy: Trigger the deploy with existing image set on App Runner
              all: Push image and trigger deploy

  -h --help   This help message.

EOF
}

log() {
  echo "$@" >&2
}

get_terraform_output() {
  id="$1"
  pushd "$MY_DIR/../terraform" > /dev/null
  log "Obtaining '$id' from terraform output ..."
  ecr_url=$(terraform output -raw "$id")
  popd > /dev/null
  echo "$ecr_url"
}

push_image() {
  ecr_url=$(get_terraform_output "ecr_url")
  region=${ecr_url#*ecr.}
  region=${region%.amazonaws.com*}

  log "Pushing image to $ecr_url ..."
  aws ecr get-login-password --region "$region" | docker login --username AWS --password-stdin "$ecr_url"

  pushd "$MY_DIR/.." > /dev/null
  docker build . -t bite-tracker
  docker tag bite-tracker "$ecr_url":latest
  docker push "$ecr_url":latest
  popd > /dev/null
}

trigger_deploy() {
  apprunner_arn=$(get_terraform_output "apprunner_arn")
  log "Deploying to App Runner '$apprunner_arn' ..."
  aws apprunner start-deployment --service-arn "$apprunner_arn"
}

deploy() {
  target="$1"

  case "$target" in
    push-image)
      push_image
      ;;
    trigger-deploy)
      trigger_deploy
      ;;
    all)
      push_image
      trigger_deploy
      ;;
    *)
      log "ERROR: Unknown target: '$target'"
      usage
      exit 1
      ;;
  esac
}

main() {
  if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    usage
    exit 0
  fi

  if [ $# -ne 1 ]; then
    log "ERROR: Invalid arguments" >&2
    usage
    exit 1
  fi

  deploy "$@"
}

main "$@"

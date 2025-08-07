#!/usr/bin/env bash
set -e

MY_DIR="$(dirname "$0")"
MY_NAME=$(basename "$0")

usage() {
  cat <<EOF
Usage: ${MY_NAME}

Deploy demo_app to AWS


OPTIONS:
  -h --help          This help message.

EOF
}

deploy() {
  echo "Obtaining arns for deployment ..."
  apprunner_arn=$(terraform output -raw apprunner_arn)
  ecr_url=$(terraform output -raw ecr_url)
  echo "Deploying to:"
  echo "App Runner.: $apprunner_arn"
  echo "ECR........: $ecr_url"
  echo

  cd "$MY_DIR"/../demo_app
  docker build . -t bite-tracker
  docker tag bite-tracker "$ecr_url":latest
  docker push "$ecr_url":latest

  echo "Deploying to App Runner '$apprunner_arn' ..."
  aws apprunner start-deployment --service-arn "$apprunner_arn"
}

main() {
  if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    usage
    exit 0
  fi

  if [ -n "$1" ]; then
    echo "ERROR: Invalid arguments"
    usage
    exit 1
  fi

  deploy
}

main "$@"

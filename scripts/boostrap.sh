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

main() {
  if [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
    usage
    exit 0
  fi

  pushd "$MY_DIR/../terraform" > /dev/null
  [ -f terraform.tfvars ] || die "terraform.tfvars file is missing"

  echo "Bootstrapping the whole environment"
  echo "-----------------------------------"

  echo "Applying ECR and RDS first ..."
  terraform init
  terraform apply \
    -target aws_ecr_repository.bite_tracker \
    -target aws_db_instance.bite_tracker \
    -target aws_instance.jumpbox \
    -target=aws_iam_role_policy_attachment.jumpbox_ssm

  echo
  exit

  echo "Bootstrapping DB roles ..."
  db_admin_password=$(grep db_admin_password terraform.tfvars | cut -f2 -d\")
  db_app_user_password=$(grep db_app_user_password terraform.tfvars | cut -f2 -d\")
  db_migrate_user_password=$(grep db_migrate_user_password terraform.tfvars | cut -f2 -d\")

  # TODO
  # docker exec -i bite-tracker-db-1 psql -U biteapp -d bite_tracker -v bt_app_user_password="'test123'" -v bt_migrate_user_password="'test1234'" < scripts/bootstrap.sql
  popd > /dev/null
  # aws ssm start-session \
  # --target i-xxxxxxxxxxxxxxx \
  # --document-name AWS-StartPortForwardingSessionToRemoteHost \
  # --parameters '{"host":["<RDS-ENDPOINT>"],"portNumber":["5432"],"localPortNumber":["5432"]}'
  #
  aws ssm send-command \
  --instance-id "i-02a05d853d64068c5" \
  --document-name "AWS-RunShellScript" \
  --comment "Run psql bootstrap" \
  --parameters 'commands=["pws && ls"]'




  echo "Pushing image to ECR ..."
  "$MY_DIR"/deploy.sh push-image
  echo

  echo "Applying rest of infra ..."
  pushd "$MY_DIR/../terraform" > /dev/null
  terraform apply
  popd > /dev/null
}

main "$@"

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

read_secret() {
  id="$1"
  aws secretsmanager get-secret-value \
      --secret-id "$id" \
      --query SecretString --output text
}

requirements() {
  type -f terraform > /dev/null || die "terraform is not installed"
  type -f docker > /dev/null || die "docker is not installed"
  type -f aws > /dev/null || die "aws cli is not installed"
  type -f jq > /dev/null || die "jq is not installed"
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
    # -target aws_db_instance.bite_tracker
    # -target aws_instance.jumpbox \
    # -target aws_iam_role_policy_attachment.jumpbox_ssm \
		# -target aws_secretsmanager_secret_version.db_admin_user_password \
		# -target aws_secretsmanager_secret_version.db_app_user_password \
		# -target aws_secretsmanager_secret_version.db_migrate_user_password

  # echo "Bootstrapping DB roles..."
  # echo "--------------------------"
  # echo "Retrieving jumpbox instance id ..."
  # jumpbox_instance_id=$(terraform show -json | jq -r ".values.outputs.jumpbox_instance_id.value")
  #
  #
  # echo "Reading DB passwords ..."
  # local db_admin_password=$(read_secret bt_db_admin_user_password)
  # local db_app_password=$(read_secret bt_db_app_user_password)
  # local db_migrate_password=$(read_secret bt_db_migrate_user_password)
  #
	# bootstrap_sql=$(<../scripts/bootstrap.sql)
	# echo "bootstrap sql $bootstrap_sql"
  #
  # echo "Creating DB roles ..."
  # aws ssm send-command \
  # --instance-id "$jumpbox_instance_id" \
  # --document-name "AWS-RunShellScript" \
  # --comment "Bootstrap DB roles" \
  # --parameters 'commands=["psql -U biteapp -d bite_tracker "]'
  #
  # aws ssm list-command-invocations \
  #   --command-id "63447ff4-e597-4cd9-9c0c-3bd0926c955e" \
  #   --details
  #

  # TODO
  # docker exec -i bite-tracker-db-1 psql -U biteapp -d bite_tracker -v bt_app_user_password="'test123'" -v bt_migrate_user_password="'test1234'" < scripts/bootstrap.sql
  popd > /dev/null
  # aws ssm start-session \
  # --target i-xxxxxxxxxxxxxxx \
  # --document-name AWS-StartPortForwardingSessionToRemoteHost \
  # --parameters '{"host":["<RDS-ENDPOINT>"],"portNumber":["5432"],"localPortNumber":["5432"]}'
  #
  # aws ssm send-command \
  # --instance-id "i-02a05d853d64068c5" \
  # --document-name "AWS-RunShellScript" \
  # --comment "Run psql bootstrap" \
  # --parameters 'commands=["pws && ls"]'


  echo "Pushing image to ECR ..."
  "$MY_DIR"/deploy.sh push-image
  echo

  echo "Applying rest of infra ..."
  pushd "$MY_DIR/../terraform" > /dev/null
	# Make sure initial db roles are creaated
  terraform apply -var bootstrap_db_roles=true
  popd > /dev/null
}

main "$@"

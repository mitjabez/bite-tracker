resource "aws_iam_role_policy_attachment" "bite_tracker_access_role_policy_attachment" {
  role       = aws_iam_role.bite_tracker_access_role.name
  policy_arn = aws_iam_policy.bite_tracker_access_policy.arn
}

resource "aws_iam_role_policy_attachment" "bite_tracker_instance_role_policy_attachment" {
  role       = aws_iam_role.bite_tracker_instance_role.name
  policy_arn = aws_iam_policy.bite_tracker_instance_policy.arn
}

resource "aws_iam_role" "bite_tracker_access_role" {
  name = "${local.name}-apprunner-access-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service : [
            "build.apprunner.amazonaws.com"
          ]
        }
      },
    ]
  })
  tags = local.default_tags
}

resource "aws_iam_role" "bite_tracker_instance_role" {
  name = "${local.name}-apprunner-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service : [
            "tasks.apprunner.amazonaws.com"
          ]
        }
      },
    ]
  })
  tags = local.default_tags
}

resource "aws_iam_policy" "bite_tracker_access_policy" {
  name = "${local.name}-apprunner-access-policy"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ecr:GetAuthorizationToken"
        ],
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = [
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:DescribeImages",
        ],
        Resource = aws_ecr_repository.bite_tracker.arn
      },
    ]
  })
  tags = local.default_tags
}

resource "aws_iam_policy" "bite_tracker_instance_policy" {
  name = "${local.name}-apprunner-instance-policy"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "secretsmanager:GetSecretValue",
          "kms:Decrypt*",
        ],
        Resource = [
          aws_secretsmanager_secret.db_app_user_username.arn,
          aws_secretsmanager_secret.db_app_user_password.arn,
          aws_secretsmanager_secret.db_admin_user_username.arn,
          aws_secretsmanager_secret.db_admin_user_password.arn,
          aws_secretsmanager_secret.hmac_token_secret.arn,
        ]
      }
    ]
  })
  tags = local.default_tags
}

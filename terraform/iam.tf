resource "aws_iam_role_policy_attachment" "bite_tracker" {
  role       = aws_iam_role.bite_tracker.name
  policy_arn = aws_iam_policy.bite_tracker.arn
}

resource "aws_iam_role" "bite_tracker" {
  name = "bite-tracker"

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

resource "aws_iam_policy" "bite_tracker" {
  name = "bite-tracker"
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
        Resource = "${aws_ecr_repository.arn}/*"
      }
    ]
  })
  tags = local.default_tags
}

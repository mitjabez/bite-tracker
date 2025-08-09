output "apprunner_service_url" {
  value = aws_apprunner_service.bite_tracker.service_url
}

output "apprunner_arn" {
  value = aws_apprunner_service.bite_tracker.arn
}

output "ecr_url" {
  value = aws_ecr_repository.bite_tracker.repository_url
}


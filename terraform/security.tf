resource "aws_security_group" "bite_tracker_db" {
  name        = "${local.name}-db"
  description = "Allow PostgreSQL from Bite Tracker"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Postgres from Bite Tracker Cluster"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [module.bite_tracker_eks.node_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.default_tags
}

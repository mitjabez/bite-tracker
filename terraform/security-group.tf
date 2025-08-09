resource "aws_security_group" "bite_tracker_app" {
  name        = "bite_tracker_app"
  description = "Allow all egress from Bite Tracker app"

  vpc_id = aws_vpc.main.id
  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = local.default_tags
}

resource "aws_security_group" "jumpbox" {
  name        = "jumpbox"
  description = "Allow all egress from jumpbox"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "Allow SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = local.default_tags
}

resource "aws_security_group" "bite_tracker_db" {
  name        = "bite-tracker-db"
  description = "Allow PostgreSQL from Bite Tracker"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Postgres from Bite Tracker"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.bite_tracker_app.id, aws_security_group.jumpbox.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.default_tags
}

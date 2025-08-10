resource "aws_db_instance" "bite_tracker" {
  allocated_storage      = 20
  identifier             = local.name
  db_name                = replace(local.name, "-", "_")
  engine                 = "postgres"
  engine_version         = "17.4"
  instance_class         = "db.t4g.micro"
  username               = var.db_admin_username
  password               = var.db_admin_password
  multi_az               = false
  availability_zone      = local.zone_b
  vpc_security_group_ids = [aws_security_group.bite_tracker_db.id]

  publicly_accessible  = false
  db_subnet_group_name = aws_db_subnet_group.bite_tracker.name
  skip_final_snapshot  = true
  apply_immediately    = true
}

resource "aws_db_subnet_group" "bite_tracker" {
  name       = local.name
  subnet_ids = [aws_subnet.bite_tracker_db_a.id, aws_subnet.bite_tracker_db_b.id]

  tags = local.default_tags
}

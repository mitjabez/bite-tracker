resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"

  # TODO: Change name
  tags = {
    Name = "main"
  }
}

resource "aws_subnet" "bite_tracker_db_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.10.0/28"
  map_public_ip_on_launch = false
  availability_zone       = local.zone_a
  tags                    = merge({ Name = "${local.name}-a-db" }, local.default_tags)
}

resource "aws_subnet" "bite_tracker_db_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.10.16/28"
  map_public_ip_on_launch = false
  availability_zone       = local.zone_b
  tags                    = merge({ Name = "${local.name}-b-db" }, local.default_tags)
}

resource "aws_subnet" "bite_tracker_app_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.11.0/28"
  map_public_ip_on_launch = false
  availability_zone       = local.zone_a
  tags                    = merge({ Name = "${local.name}-a-app" }, local.default_tags)
}


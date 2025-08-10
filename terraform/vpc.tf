resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"

  tags = {
    Name = "main"
  }
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/28"
  map_public_ip_on_launch = false
  availability_zone       = local.zone_a
  tags                    = merge({ Name = "public-a" }, local.default_tags)
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


# TODO: Cleanup SSM stuff. Do we event need vpc gw?
locals {
  services = {
    "ec2messages" : {
      "name" : "com.amazonaws.${local.region}.ec2messages"
    },
    "ssm" : {
      "name" : "com.amazonaws.${local.region}.ssm"
    },
    "ssmmessages" : {
      "name" : "com.amazonaws.${local.region}.ssmmessages"
    }
  }
}

resource "aws_vpc_endpoint" "ssm_endpoint" {
  for_each           = local.services
  vpc_id             = aws_vpc.main.id
  service_name       = each.value.name
  vpc_endpoint_type  = "Interface"
  security_group_ids = [aws_security_group.ssm_https.id]
  ip_address_type    = "ipv4"
  subnet_ids         = [aws_subnet.public.id]
}

resource "aws_security_group" "ssm_https" {
  name        = "allow_ssm"
  description = "Allow SSM traffic"
  vpc_id      = aws_vpc.main.id
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_subnet.public.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = local.default_tags
}

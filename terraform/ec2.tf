resource "aws_instance" "jumpbox" {
  # Amazon Linux 2023 AMI 2023.8.20250808.1 x86_64 HVM kernel-6.1
  ami                         = "ami-05a2d2d0a1020fecd"
  instance_type               = "t3.micro"
  subnet_id                   = aws_subnet.public.id
  vpc_security_group_ids      = [aws_security_group.jumpbox.id]
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.jumpbox.name
  key_name                    = aws_key_pair.jumpbox_admin.key_name

  user_data = <<-EOF
    #!/usr/bin/env bash
    sudo dnf update -y
    sudo dnf install -y postgresql15
  EOF

  tags = merge({ Name = "jumpbox" }, local.default_tags)
}

# TODO: Remove
resource "aws_key_pair" "jumpbox_admin" {
  key_name   = "jumpbox-admin"
  public_key = file(var.ssh_public_key_path)
}

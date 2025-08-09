data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "jumpbox" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t3.micro"
  subnet_id                   = aws_subnet.dmz.id
  security_groups             = [aws_security_group.jumpbox.id]
  associate_public_ip_address = true
  key_name                    = aws_key_pair.jumpbox_admin.key_name
  tags                        = merge({ Name = "jumpbox" }, local.default_tags)
}

resource "aws_key_pair" "jumpbox_admin" {
  key_name   = "jumpbox-admin"
  public_key = file(var.ssh_public_key_path)
}

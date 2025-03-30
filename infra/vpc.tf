resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  enable_dns_support = true
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_subnet" "public_a" {
  vpc_id                    = aws_vpc.main.id
  cidr_block                = "10.0.1.0/24"
  availability_zone         = "us-west-2a"
  map_public_ip_on_launch   = true
}

resource "aws_subnet" "public_b" {
  vpc_id                    = aws_vpc.main.id
  cidr_block                = "10.0.2.0/24"
  availability_zone         = "us-west-2b"
  map_public_ip_on_launch   = true
}

resource "aws_security_group" "rds" {
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # 👈 only for testing! lock this down in prod
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
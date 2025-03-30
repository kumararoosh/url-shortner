resource "aws_db_subnet_group" "shortener" {
  name       = "shortener-db-subnet-group"
  subnet_ids = [aws_subnet.public_a.id, aws_subnet.public_b.id]
}

resource "aws_db_instance" "shortener" {
  identifier         = "shortener-db"
  allocated_storage  = 20
  engine             = "postgres"
  engine_version     = "17.2"
  instance_class     = "db.t3.micro"
  username           = var.db_username
  password           = var.db_password
  db_name            = var.db_name
  publicly_accessible = true
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.shortener.name
  skip_final_snapshot    = true
}

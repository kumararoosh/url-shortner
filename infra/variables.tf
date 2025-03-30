variable "aws_region" {
    default = "us-west-2"
}

variable "db_username" {
    default = "shorty"
}

variable "db_password" {
    description = "RDS database password"
    sensitive = true
}

variable "db_name" {
  default = "shortener"
}
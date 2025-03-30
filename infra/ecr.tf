resource "aws_ecr_repository" "shortener" {
  name = "url-shortener"
  image_scanning_configuration {
    scan_on_push = true
  }
}
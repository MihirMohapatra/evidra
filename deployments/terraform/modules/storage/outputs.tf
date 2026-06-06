output "endpoint" {
  value = "https://s3.${data.aws_region.current.name}.amazonaws.com"
}

output "bucket_names" {
  value = values(aws_s3_bucket.this)[*].id
}

output "bucket_arns" {
  value = values(aws_s3_bucket.this)[*].arn
}

data "aws_region" "current" {}

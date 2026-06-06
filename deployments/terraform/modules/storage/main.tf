resource "aws_s3_bucket" "this" {
  for_each = toset(var.buckets)

  bucket = "${var.environment}-evidra-${each.key}"

  force_destroy = var.force_destroy

  tags = merge(var.tags, {
    Name        = "${var.environment}-evidra-${each.key}"
    Environment = var.environment
    Bucket      = each.key
  })
}

resource "aws_s3_bucket_versioning" "this" {
  for_each = var.enable_versioning ? aws_s3_bucket.this : {}

  bucket = each.value.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  for_each = var.enable_encryption ? aws_s3_bucket.this : {}

  bucket = each.value.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = var.kms_key_arn != "" ? "aws:kms" : "AES256"
      kms_master_key_id = var.kms_key_arn != "" ? var.kms_key_arn : null
    }
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  for_each = aws_s3_bucket.this

  bucket = each.value.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  for_each = aws_s3_bucket.this

  bucket = each.value.id

  rule {
    id     = "abort-incomplete-multipart"
    status = "Enabled"

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
}

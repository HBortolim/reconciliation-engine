# CloudFront Distribution for dashboard
resource "aws_cloudfront_distribution" "dashboard" {
  enabled = true

  origin {
    domain_name = "${var.storage_bucket_name}.s3.${var.aws_region}.amazonaws.com"
    origin_id   = "s3-dashboard"

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.dashboard.cloudfront_access_identity_path
    }
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "s3-dashboard"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = var.price_class

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = {
    Name = "${var.project_name}-cdn"
  }
}

# CloudFront Origin Access Identity
resource "aws_cloudfront_origin_access_identity" "dashboard" {
  comment = "OAI for ${var.project_name} dashboard"
}

# S3 Bucket Policy for CloudFront
resource "aws_s3_bucket_policy" "cloudfront_access" {
  bucket = var.storage_bucket_name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "CloudFrontAccess"
        Effect = "Allow"
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.dashboard.iam_arn
        }
        Action   = "s3:GetObject"
        Resource = "${var.storage_bucket_arn}/*"
      }
    ]
  })
}

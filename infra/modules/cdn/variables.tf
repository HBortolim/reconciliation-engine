variable "project_name" {
  description = "Project name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "storage_bucket_name" {
  description = "S3 bucket name"
  type        = string
}

variable "storage_bucket_arn" {
  description = "S3 bucket ARN"
  type        = string
}

variable "price_class" {
  description = "CloudFront price class (PriceClass_All, PriceClass_100, PriceClass_200)"
  type        = string
  default     = "PriceClass_100"
}

variable "project_name" {
  description = "Project name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs"
  type        = list(string)
}

variable "alb_security_group_id" {
  description = "Security group ID for ALB"
  type        = string
}

variable "ecs_security_group_id" {
  description = "Security group ID for ECS tasks"
  type        = string
}

variable "api_image" {
  description = "Docker image for API"
  type        = string
}

variable "worker_image" {
  description = "Docker image for Worker"
  type        = string
}

variable "api_task_cpu" {
  description = "CPU units for API task"
  type        = string
  default     = "256"
}

variable "api_task_memory" {
  description = "Memory (MB) for API task"
  type        = string
  default     = "512"
}

variable "worker_task_cpu" {
  description = "CPU units for Worker task"
  type        = string
  default     = "256"
}

variable "worker_task_memory" {
  description = "Memory (MB) for Worker task"
  type        = string
  default     = "512"
}

variable "api_desired_count" {
  description = "Desired number of API tasks"
  type        = number
  default     = 1
}

variable "worker_desired_count" {
  description = "Desired number of Worker tasks"
  type        = number
  default     = 1
}

variable "database_host" {
  description = "Database host"
  type        = string
}

variable "database_port" {
  description = "Database port"
  type        = string
}

variable "database_name" {
  description = "Database name"
  type        = string
}

variable "database_password_secret_arn" {
  description = "ARN of Secrets Manager secret for database password"
  type        = string
}

variable "redis_host" {
  description = "Redis host"
  type        = string
}

variable "redis_port" {
  description = "Redis port"
  type        = string
}

variable "storage_bucket_name" {
  description = "S3 bucket name for storage"
  type        = string
}

variable "storage_bucket_arn" {
  description = "S3 bucket ARN for storage"
  type        = string
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 7
}

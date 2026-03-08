variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = "reconciliation-engine"
}

# Networking
variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.11.0/24"]
}

# Database
variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.r6g.xlarge"
}

variable "db_allocated_storage" {
  description = "Allocated storage in GB"
  type        = number
  default     = 100
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "reconciliation_db"
}

variable "db_username" {
  description = "Database master username"
  type        = string
  default     = "postgres"
  sensitive   = true
}

variable "db_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

variable "db_backup_retention_days" {
  description = "Backup retention days"
  type        = number
  default     = 30
}

variable "db_multi_az" {
  description = "Enable Multi-AZ"
  type        = bool
  default     = true
}

variable "db_deletion_protection" {
  description = "Enable deletion protection"
  type        = bool
  default     = true
}

# Cache
variable "redis_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.r7g.large"
}

variable "redis_num_nodes" {
  description = "Number of Redis nodes"
  type        = number
  default     = 3
}

variable "redis_automatic_failover" {
  description = "Enable automatic failover"
  type        = bool
  default     = true
}

variable "redis_multi_az" {
  description = "Enable Multi-AZ"
  type        = bool
  default     = true
}

# Compute
variable "api_image" {
  description = "Docker image for API"
  type        = string
  default     = "reconciliation-engine-api:latest"
}

variable "worker_image" {
  description = "Docker image for Worker"
  type        = string
  default     = "reconciliation-engine-worker:latest"
}

variable "api_task_cpu" {
  description = "API task CPU units"
  type        = string
  default     = "1024"
}

variable "api_task_memory" {
  description = "API task memory (MB)"
  type        = string
  default     = "2048"
}

variable "worker_task_cpu" {
  description = "Worker task CPU units"
  type        = string
  default     = "512"
}

variable "worker_task_memory" {
  description = "Worker task memory (MB)"
  type        = string
  default     = "1024"
}

variable "api_desired_count" {
  description = "Desired number of API tasks"
  type        = number
  default     = 3
}

variable "worker_desired_count" {
  description = "Desired number of Worker tasks"
  type        = number
  default     = 5
}

# CDN
variable "cloudfront_price_class" {
  description = "CloudFront price class"
  type        = string
  default     = "PriceClass_All"
}

# Observability
variable "log_retention_days" {
  description = "CloudWatch log retention days"
  type        = number
  default     = 30
}

variable "alarm_email" {
  description = "Email address for CloudWatch alarms"
  type        = string
  default     = null
}

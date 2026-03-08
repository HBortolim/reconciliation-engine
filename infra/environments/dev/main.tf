terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }

  backend "s3" {
    bucket         = "reconciliation-engine-terraform-state"
    key            = "dev/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "reconciliation-engine-terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "reconciliation-engine"
      Environment = "dev"
      ManagedBy   = "Terraform"
    }
  }
}

# Networking Module
module "networking" {
  source = "../../modules/networking"

  project_name           = var.project_name
  vpc_cidr               = var.vpc_cidr
  public_subnet_cidrs    = var.public_subnet_cidrs
  private_subnet_cidrs   = var.private_subnet_cidrs
}

# Database Module
module "database" {
  source = "../../modules/database"

  project_name           = var.project_name
  instance_class         = var.db_instance_class
  allocated_storage      = var.db_allocated_storage
  db_name                = var.db_name
  db_username            = var.db_username
  db_password            = var.db_password
  subnet_ids             = module.networking.private_subnet_ids
  rds_security_group_id  = module.networking.rds_security_group_id
  backup_retention_days  = var.db_backup_retention_days
  multi_az               = var.db_multi_az
  deletion_protection    = var.db_deletion_protection
}

# Cache Module
module "cache" {
  source = "../../modules/cache"

  project_name                 = var.project_name
  node_type                    = var.redis_node_type
  num_cache_nodes              = var.redis_num_nodes
  subnet_ids                   = module.networking.private_subnet_ids
  elasticache_security_group_id = module.networking.elasticache_security_group_id
  automatic_failover_enabled   = var.redis_automatic_failover
  multi_az_enabled             = var.redis_multi_az
}

# Storage Module
module "storage" {
  source = "../../modules/storage"

  project_name = var.project_name
}

# Compute Module
module "compute" {
  source = "../../modules/compute"

  project_name                  = var.project_name
  aws_region                    = var.aws_region
  vpc_id                        = module.networking.vpc_id
  public_subnet_ids             = module.networking.public_subnet_ids
  private_subnet_ids            = module.networking.private_subnet_ids
  alb_security_group_id         = module.networking.alb_security_group_id
  ecs_security_group_id         = module.networking.ecs_security_group_id
  api_image                     = var.api_image
  worker_image                  = var.worker_image
  api_task_cpu                  = var.api_task_cpu
  api_task_memory               = var.api_task_memory
  worker_task_cpu               = var.worker_task_cpu
  worker_task_memory            = var.worker_task_memory
  api_desired_count             = var.api_desired_count
  worker_desired_count          = var.worker_desired_count
  database_host                 = module.database.address
  database_port                 = tostring(module.database.port)
  database_name                 = module.database.db_name
  database_password_secret_arn  = aws_secretsmanager_secret.db_password.arn
  redis_host                    = module.cache.endpoint
  redis_port                    = tostring(module.cache.port)
  storage_bucket_name           = module.storage.bucket_name
  storage_bucket_arn            = module.storage.bucket_arn
}

# CDN Module
module "cdn" {
  source = "../../modules/cdn"

  project_name           = var.project_name
  aws_region             = var.aws_region
  storage_bucket_name    = module.storage.bucket_name
  storage_bucket_arn     = module.storage.bucket_arn
  price_class            = var.cloudfront_price_class
}

# Observability Module
module "observability" {
  source = "../../modules/observability"

  project_name              = var.project_name
  aws_region                = var.aws_region
  log_retention_days        = var.log_retention_days
  alarm_email               = var.alarm_email
  db_instance_identifier    = split(":", module.database.endpoint)[0]
  ecs_cluster_name          = module.compute.ecs_cluster_name
  ecs_service_name          = module.compute.api_service_name
  target_group_name         = "reconciliation-engine-api-tg"
  load_balancer_name        = "reconciliation-engine-alb"
}

# Secrets Manager Secret for Database Password
resource "aws_secretsmanager_secret" "db_password" {
  name                    = "${var.project_name}-db-password"
  recovery_window_in_days = 7

  tags = {
    Name = "${var.project_name}-db-password"
  }
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id       = aws_secretsmanager_secret.db_password.id
  secret_string   = var.db_password
}

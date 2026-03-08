aws_region = "us-east-1"
project_name = "reconciliation-engine"

# Networking
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.11.0/24"]

# Database - Dev: minimal instance
db_instance_class        = "db.t4g.micro"
db_allocated_storage     = 20
db_name                  = "reconciliation_db"
db_backup_retention_days = 7
db_multi_az              = false
db_deletion_protection   = false

# Cache - Dev: single node
redis_node_type           = "cache.t4g.micro"
redis_num_nodes           = 1
redis_automatic_failover  = false
redis_multi_az            = false

# Compute - Dev: minimal tasks
api_task_cpu       = "256"
api_task_memory    = "512"
worker_task_cpu    = "256"
worker_task_memory = "512"
api_desired_count  = 1
worker_desired_count = 1

# Docker Images (update with actual ECR URIs)
api_image    = "reconciliation-engine-api:latest"
worker_image = "reconciliation-engine-worker:latest"

# CDN
cloudfront_price_class = "PriceClass_100"

# Observability
log_retention_days = 7

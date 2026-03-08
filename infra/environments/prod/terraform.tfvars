aws_region = "us-east-1"
project_name = "reconciliation-engine"

# Networking
vpc_cidr             = "10.0.0.0/16"
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.10.0/24", "10.0.11.0/24"]

# Database - Prod: high-performance instance
db_instance_class        = "db.r6g.xlarge"
db_allocated_storage     = 100
db_name                  = "reconciliation_db"
db_backup_retention_days = 30
db_multi_az              = true
db_deletion_protection   = true

# Cache - Prod: multi-node cluster with failover
redis_node_type           = "cache.r7g.large"
redis_num_nodes           = 3
redis_automatic_failover  = true
redis_multi_az            = true

# Compute - Prod: scaled instances
api_task_cpu       = "1024"
api_task_memory    = "2048"
worker_task_cpu    = "512"
worker_task_memory = "1024"
api_desired_count    = 3
worker_desired_count = 5

# Docker Images (update with actual ECR URIs)
api_image    = "reconciliation-engine-api:latest"
worker_image = "reconciliation-engine-worker:latest"

# CDN
cloudfront_price_class = "PriceClass_All"

# Observability
log_retention_days = 30
# alarm_email = "ops-team@example.com"  # Uncomment and set to your alert email

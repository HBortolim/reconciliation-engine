output "endpoint" {
  description = "ElastiCache cluster endpoint"
  value       = aws_elasticache_cluster.main.cache_nodes[0].address
}

output "port" {
  description = "ElastiCache cluster port"
  value       = aws_elasticache_cluster.main.port
}

output "cluster_id" {
  description = "ElastiCache cluster ID"
  value       = aws_elasticache_cluster.main.cluster_id
}

output "engine_version" {
  description = "ElastiCache engine version"
  value       = aws_elasticache_cluster.main.engine_version
}

output "sns_topic_arn" {
  description = "SNS topic ARN for alarms"
  value       = aws_sns_topic.alarms.arn
}

output "log_group_application" {
  description = "CloudWatch log group for application"
  value       = aws_cloudwatch_log_group.application.name
}

output "log_group_database" {
  description = "CloudWatch log group for database"
  value       = aws_cloudwatch_log_group.database.name
}

output "dashboard_url" {
  description = "CloudWatch dashboard URL"
  value       = "https://console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${aws_cloudwatch_dashboard.main.dashboard_name}"
}

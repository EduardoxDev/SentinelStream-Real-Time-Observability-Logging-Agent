output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "ecs_cluster_name" {
  description = "ECS Cluster name"
  value       = aws_ecs_cluster.main.name
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = aws_elasticache_cluster.redis.cache_nodes[0].address
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = aws_lb.server.dns_name
}

output "cloudwatch_log_group_agent" {
  description = "CloudWatch Log Group for Agent"
  value       = aws_cloudwatch_log_group.agent.name
}

output "cloudwatch_log_group_server" {
  description = "CloudWatch Log Group for Server"
  value       = aws_cloudwatch_log_group.server.name
}

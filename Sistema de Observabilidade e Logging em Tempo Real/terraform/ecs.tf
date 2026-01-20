# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "observability-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Environment = var.environment
  }
}

# ECS Task Definition - Agent
resource "aws_ecs_task_definition" "agent" {
  family                   = "observability-agent"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name  = "agent"
      image = "${var.ecr_repository_url}:agent-latest"
      
      environment = [
        {
          name  = "INFLUXDB_URL"
          value = "http://${aws_lb.influxdb.dns_name}:8086"
        },
        {
          name  = "REDIS_ADDR"
          value = "${aws_elasticache_cluster.redis.cache_nodes[0].address}:6379"
        }
      ]

      secrets = [
        {
          name      = "INFLUXDB_TOKEN"
          valueFrom = aws_secretsmanager_secret.influxdb_token.arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.agent.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "agent"
        }
      }

      mountPoints = [
        {
          sourceVolume  = "docker_sock"
          containerPath = "/var/run/docker.sock"
          readOnly      = true
        }
      ]
    }
  ])

  volume {
    name      = "docker_sock"
    host_path = "/var/run/docker.sock"
  }
}

# ECS Task Definition - Server
resource "aws_ecs_task_definition" "server" {
  family                   = "observability-server"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024"
  memory                   = "2048"
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name  = "server"
      image = "${var.ecr_repository_url}:server-latest"
      
      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "INFLUXDB_URL"
          value = "http://${aws_lb.influxdb.dns_name}:8086"
        },
        {
          name  = "REDIS_ADDR"
          value = "${aws_elasticache_cluster.redis.cache_nodes[0].address}:6379"
        }
      ]

      secrets = [
        {
          name      = "INFLUXDB_TOKEN"
          valueFrom = aws_secretsmanager_secret.influxdb_token.arn
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.server.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "server"
        }
      }
    }
  ])
}

# ECS Service - Server
resource "aws_ecs_service" "server" {
  name            = "observability-server"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.server.arn
  desired_count   = 2
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private[*].id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.server.arn
    container_name   = "server"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.server]
}

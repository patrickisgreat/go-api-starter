module "api_service" {
  # TODO: Replace with alternative Terraform module - was: source = "github.com/soundcloud/tf-common//modules/aws-ecs-service?ref=aws-ecs-service@3.48.0"
  source       = "github.com/soundcloud/tf-common//modules/aws-ecs-service?ref=aws-ecs-service@3.48.0"
  service_name = var.system
  aws_region   = local.aws_region
  aws_profile  = var.aws_profile
  cluster_name = local.ecs_cluster
  environment  = local.environment
  vpc_id       = data.aws_vpc.sc.id
  cpu          = var.cpu
  memory       = var.memory
  subnets      = data.aws_subnets.private.ids
  environment_variables = merge(local.telemetry_config, {
    APP_NAME    = var.system
    LOG_LEVEL   = local.log_level
    ENVIRONMENT = local.environment
    PORT        = var.container_port
  })
  docker_image_url   = local.container_image
  ecs_min_tasks      = var.ecs_min_tasks
  ecs_desired_tasks  = var.ecs_desired_tasks
  ecs_max_tasks      = var.ecs_max_tasks
  scaleup_cooldown   = var.scaleup_cooldown
  scaledown_cooldown = var.scaledown_cooldown
  capacity_provider_strategy = {
    FARGATE      = local.on_demand_capacity
    FARGATE_SPOT = local.spot_capacity
  }
  minimum_non_spot_tasks = 1
  cpu_architecture       = "ARM64"
  target_tracking_policy = {
    target_value           = 40
    scaleup_cooldown       = 60
    scaledown_cooldown     = 60
    predefined_metric_type = "ECSServiceAverageCPUUtilization"
    resource_label         = null
  }
  target_tracking_policy_enabled = true
  use_docker_healthcheck         = false
  virtual_node_healthcheck_config = {
    healthy_threshold   = 3
    unhealthy_threshold = 2
    interval            = 5
    timeout             = 2
    path                = var.api_healthcheck_path
  }
  enable_execute_command            = true
  enable_app_mesh                   = true
  mesh_name                         = local.short_env_name
  debug_app_mesh                    = local.environment == "development"
  ingress_cidr_ranges               = ["10.0.0.0/8"]
  mesh_backends                     = []
  policies                          = []
  container_port                    = var.container_port
  envoy_listener_max_retries        = 5
  envoy_image_version               = "v1.29.6.1-prod"
  container_secrets_manager_secrets = []
  enable_new_relic_sidecar          = false
  enable_otel_collector_sidecar     = true
  enable_ecs_container_metrics      = false
  tags = merge(local.tags, {
    Component = "api"
  })
}

module "alerts" {
  # TODO: Replace with alternative Terraform module - was: source = "github.com/soundcloud/tf-common//modules/grafana-alerts?ref=grafana-alerts@1.5.0"
  source = "github.com/soundcloud/tf-common//modules/grafana-alerts?ref=grafana-alerts@1.5.0"

  count        = var.alerts_enabled ? 1 : 0
  system       = var.system
  component    = "api"
  environment  = var.environment
  service_name = "${var.system}-api"

  owner = "dx-team"
  repo  = "patrickisgreat/go-api-starter"

  slack_channel          = "dx-alerts"
  enable_pagerduty       = false
  pagerduty_service_name = "none"

  alert_on_no_data = false
}

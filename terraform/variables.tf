variable "aws_profile" {
  description = "AWS account profile"
  type        = string
}

variable "commit_sha" {
  description = "The commit being deployed, used for the docker image tags"
  type        = string
}

variable "environment" {
  description = "The environment to deploy to"
  type        = string
}

# ECS
variable "system" {
  description = "The name of the system"
  type        = string
  default     = "go-api-starter"
}

variable "cpu" {
  description = "Task CPU limit in units"
  type        = number
  default     = 1024
}

variable "memory" {
  description = "Task memory limit in MB"
  type        = number
  default     = 2048
}

variable "api_healthcheck_path" {
  description = "The route that the lambda should healthcheck on"
  type        = string
  default     = "/-/health"
}

variable "ecs_min_tasks" {
  description = "Minimum number of api containers to run per region"
  type        = number
  default     = 1
}

variable "ecs_desired_tasks" {
  description = "Desired number of api containers to run per region"
  type        = number
  default     = 1
}

variable "ecs_max_tasks" {
  description = "Maximum number of api containers to run per region"
  type        = number
  default     = 3
}

variable "scaleup_cooldown" {
  description = "Time to wait after scale up before scaling up again"
  type        = number
  default     = 60
}

variable "scaledown_cooldown" {
  description = "Time to wait after scale down before scaling down again"
  type        = number
  default     = 300
}

variable "container_port" {
  description = "Container port"
  type        = number
  default     = 8000
}

# Alerts
variable "alerts_enabled" {
  description = "Whether or not to enable alerts for the environemnt"
  type        = bool
  default     = false
}

variable "alerts_runbook_url" {
  description = "The go-api-starter runbook URL"
  type        = string
  default     = ""
}

variable "error_duration" {
  description = "Set the duration in minutes for the high errors alert"
  type        = number
  default     = 2
}

variable "latency_threshold" {
  description = "Set the threshold for the high response time alert"
  type        = string
  default     = "3"
}

variable "latency_duration" {
  description = "Set the duration in minutes for the high response time alert"
  type        = number
  default     = 2
}

variable "error_threshold" {
  description = "Set the threshold percentage for the high errors alert"
  type        = string
  default     = "5"
}

variable "slack" {
  description = <<-EOF
  Optionally define slack notification behaviour here. By default, when present, this will create a notification channel
  for slack and integrates the channel into existing workflow which includes PagerDuty destination. To just send to slack only, i.e. warning/non-paging alerts, set the warn flag to true which will create a separate workflow just for slack destination.
EOF
  type = object({
    channel = optional(string)
    warn    = optional(bool, false)
  })
  default = {
    channel = "api-team-alerts"
  }
}

variable "slack_additional" {
  description = <<-EOF
  Optionally define additional slack channel where slack alerts should go.
EOF
  type = object({
    channel = optional(string)
    warn    = optional(bool, false)
  })
  default = {}
}

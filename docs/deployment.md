# Deployment

The service ships as a container image and runs on AWS ECS. Terraform in
`terraform/` describes the service, its scaling and its alerts, with one
configuration directory per environment.

## Container image

`Dockerfile` is a two stage build:

1. **build** on `golang:alpine`. Downloads modules, then compiles a static binary
   with `CGO_ENABLED=0` and `-ldflags="-w -s"` to strip debug information.
2. **final** on `scratch`. Copies in only the CA certificate bundle and the
   binary. There is no shell and no package manager in the image.

The private `gokit` module means `go mod download` needs GitHub credentials. They
are passed as a BuildKit secret so that they never land in an image layer:

```sh
docker build \
  --secret id=git_credentials,src=$HOME/.git-credentials \
  --platform linux/amd64 \
  -t pb-go-api-starter:local .

docker run --rm -p 8000:8000 \
  -e APP_NAME=pb-go-api-starter -e HOST=0.0.0.0 -e PORT=8000 \
  -e ENVIRONMENT=local -e LOG_LEVEL=debug -e OTEL_SDK_DISABLED=true \
  pb-go-api-starter:local
```

`.dockerignore` keeps docs, terraform, env files and the task cache out of the
build context.

## Image registry

Terraform expects the image at
`145582369311.dkr.ecr.eu-central-1.amazonaws.com/pb-go-api-starter:<commit sha>`.
Tag with the full commit SHA so that every deployment is traceable to a commit
and rollbacks are a matter of applying with an older `commit_sha`.

## Terraform layout

```
terraform/
  main.tf          module "api_service" (ECS service) and module "alerts"
  variables.tf     system name, sizing, scaling, alert thresholds, runbook URL
  locals.tf        per environment log level, trace sampling, telemetry env vars, image URI
  data.tf          account, VPC and subnet lookups from SSM
  terraform.tf     backend and provider configuration
  config/
    development/   config.remote (state backend) and terraform.tfvars
    staging/
    production/
```

Environment specific values are resolved two ways. `config/<env>/terraform.tfvars`
supplies the AWS profile and environment name. `locals.tf` then maps the
environment to log level, trace ratio and capacity strategy (spot instances for
non-production, on-demand only for production).

## Applying

```sh
cd terraform
terraform init -backend-config=config/development/config.remote
terraform plan  -var-file=config/development/terraform.tfvars -var commit_sha=<sha> -out plan.tfplan
terraform apply plan.tfplan
```

Replace `development` with `staging` or `production`. The state backend for each
environment is an S3 bucket with a DynamoDB lock table, named in `config.remote`.

## Alerts

`module "alerts"` creates high error rate and high response time alerts. The
thresholds and evaluation windows are variables with defaults in `variables.tf`,
and `alerts_runbook_url` should point at the runbook for this service so that the
alert links straight to it. `GetQuoteFlaky` and `GetQuoteLate` exist so these
alerts can be triggered on purpose in a non-production environment.

## Verifying a deployment

After applying, run the end to end suite against the environment:

```sh
ENV=development task test:e2e:ci
```

This uses `API_HOST` and `API_PORT` from `.env/.env.development`. Then confirm the
service is reporting to telemetry by checking the dashboard linked from the
README.

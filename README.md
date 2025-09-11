# Bite Tracker [![Go](https://github.com/mitjabez/bite-tracker/actions/workflows/ci.yml/badge.svg)](https://github.com/mitjabez/bite-tracker/actions/workflows/ci.yml)

## About

A simple food tracker website. Created for learnging myself Go on a real world example. Stack:

Frontend:

- Htmx
- Bulma
- Alpine.js

Backend:

- Go
- sqlc
- a-h/templ
- Postgres
- golang-migrate
- testcontainers-go

## Local development

### Prerequisites

- [Docker](https://www.docker.com/): Local DB
- [air](https://github.com/air-verse/air): Live reloading
- [a-h/templ](https://templ.guide/): HTML templating
- [sqlc](https://sqlc.dev/): SQL to Go code compiler
- [golang-migrate](https://github.com/golang-migrate/migrate): DB Schema migrations
- [testcontainers-go](https://testcontainers.com/): Isolated dependencies for testing

### Running Locally

The application handles database migrations automatically on startup. To run the app locally, you first need to start the Postgres database in Docker:

```bash
make db-start
```

After the [prerequisites](#prerequisites) have been met you can run the app locally.

- Run locally with live reload and access on http://localhost:3000

```sh
air
```

- Without live reload:

```sh
make run
```

- Generate templ and sqlc files only:

```sh
make generate
```

### Testing

To run the unit tests:

```sh
make test
```

## Infrastructure

The infrastructure is defined in Terraform and is deployed to AWS. The stack consists of:

- AWS App Runner
- AWS RDS (PostgreSQL)
- AWS Secrets Manager
- AWS SSM for DB managementt

To bootstrap the infrastructure, you need to have [Terraform](https://www.terraform.io/) and
[AWS CLI](https://aws.amazon.com/cli/) installed and configured.

Before running the bootstrap script, you need to create a `terraform.tfvars` file in the `terraform` directory. You can
use `terraform.tfvars.example` as a template.

To bootstrap the infrastructure, run the following command:

```sh
./scripts/bootstrap.sh
```

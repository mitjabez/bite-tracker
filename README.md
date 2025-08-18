# Bite Tracker [![Go](https://github.com/mitjabez/bite-tracker/actions/workflows/ci.yml/badge.svg)](https://github.com/mitjabez/bite-tracker/actions/workflows/ci.yml)

## About

A simple food tracker website. Stack:

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

## Local development

### Prerequisites

- [Docker](https://www.docker.com/): Local DB
- [air](https://github.com/air-verse/air): Live reloading
- [a-h/templ](https://templ.guide/): HTML templating
- [sqlc](https://sqlc.dev/): SQL to Go code compiler
- [golang-migrate](https://github.com/golang-migrate/migrate): DB Schema migrations

### DB Init

Make sure that the database has been set before [running locally](#running-locally).

- Start Postgres DB in Docker:

```bash
make db-start
```

- Initialize schema:

```bash
make db-up
```

### Running Locally

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

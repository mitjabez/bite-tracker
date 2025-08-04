# Bite Tracker

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

## Prerequisites

- [a-h/templ](https://templ.guide/)
- [sqlc](https://sqlc.dev/)

## Usage

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

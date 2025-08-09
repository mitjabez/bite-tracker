FROM golang:1.24.4@sha256:20a022e5112a144aa7b7aeb3f22ebf2cdaefcc4aac0d64e8deeee8cdc18b9c0f AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY Makefile ./Makefile

RUN CGO_ENABLED=0 GOOS=linux make build

FROM gcr.io/distroless/static:nonroot@sha256:cdf4daaf154e3e27cfffc799c16f343a384228f38646928a1513d925f473cb46

WORKDIR /app

COPY --from=builder /app/build/bite-tracker /app/bite-tracker
COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations
EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/bite-tracker"]

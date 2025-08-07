FROM golang:1.24

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /demo_app

EXPOSE 8080

CMD ["/demo_app"]

FROM scratch
COPY --from=0 /demo_app /demo_app
CMD ["/demo_app"]

FROM golang:1.23.0 as builder
WORKDIR /app
COPY . .

RUN go mod tidy && CGO_ENABLED=0 go build -o api-gateway ./cmd/api-gateway

FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/api-gateway .

CMD ["./api-gateway"]
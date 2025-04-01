FROM golang:1.23.0 as builder
WORKDIR /app
COPY . .

# Устанавливаем зависимости и инструменты для gRPC
# RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
#     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN protoc -I api/proto --go_out=pkg/post-service-api --go-grpc_out=pkg/post-service-api api/proto/post-service.proto

RUN go mod tidy && CGO_ENABLED=0 go build -o post-service ./cmd/post-service

FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/post-service .

CMD ["./post-service"]

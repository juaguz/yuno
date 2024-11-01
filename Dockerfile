# Etapa de compilación
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o myapp ./cmd/api

FROM alpine:latest

COPY --from=builder /app/myapp /usr/local/bin/myapp

CMD ["myapp"]

# Build
FROM golang:1.23.4-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wrapto .

# Staging
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/wrapto .
COPY --from=builder /app/config/config.yml .

EXPOSE 3000

ENTRYPOINT ["./wrapto", "./config.yml"]

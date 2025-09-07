
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./finance-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./finance-consumer ./cmd/consumer/main.go

FROM alpine:latest as deploy

RUN apk update && apk add --no-cache tzdata ca-certificates && \
    addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

ENV TZ=Asia/Tehran

WORKDIR /app

COPY --from=builder /app/finance-api ./finance-api
COPY --from=builder /app/finance-consumer ./finance-consumer

COPY config.yaml ./config.yaml

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8081
CMD ["./finance-api"]

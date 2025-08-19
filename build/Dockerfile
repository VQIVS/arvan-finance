FROM golang:1.24-slim as builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./finance-api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./finance-consumer cmd/consumer/main.go

FROM alpine as deploy

RUN apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Tehran

WORKDIR /app

COPY --from=builder /app/finance-api ./finance-api
COPY --from=builder /app/finance-consumer ./finance-consumer
COPY config.json ./config.json

EXPOSE 8080

CMD ["./finance-api"]

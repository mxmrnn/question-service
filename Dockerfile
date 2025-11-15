FROM golang:1.24-alpine AS builder
WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin && \
    go build -o /app/bin/api ./cmd/api && \
    go build -o /app/bin/migrate ./cmd/migrate

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/migrate /app/migrate
COPY --from=builder /app/migrations /app/migrations

RUN chmod +x /app/api /app/migrate

EXPOSE 8080

ENV HTTP_PORT=":8080"

CMD ["/app/api"]

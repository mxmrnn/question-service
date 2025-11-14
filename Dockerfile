FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o question-service ./cmd/api


FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/question-service /app/question-service

EXPOSE 8080

ENV HTTP_PORT=":8080"

CMD ["/app/question-service"]

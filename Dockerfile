FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /weather-api ./cmd/api

FROM alpine:latest

WORKDIR /

COPY --from=builder /weather-api /weather-api

COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/web /web

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/weather-api"]
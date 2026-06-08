FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/server /app/server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
RUN echo "POSTGRES_CONNECTION_STRING=" > /app/.env && \
    echo "ACCESS_SECRET=" >> /app/.env && \
    echo "REFRESH_SECRET=" >> /app/.env && \
    echo "ENV=production" >> /app/.env
EXPOSE 8080
ENTRYPOINT ["/app/server"]

# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/rest-api-app ./cmd/api-app

# Final stage
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/bin/rest-api-app .

EXPOSE 8080

CMD ["./rest-api-app"]
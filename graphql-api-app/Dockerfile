FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod vendor
RUN go build -o /app/bin/graphql-api-app /app/cmd/api-app/

FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/bin/graphql-api-app .
EXPOSE 80
CMD ["./graphql-api-app"]`
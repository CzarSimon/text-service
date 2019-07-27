FROM golang:1.12-alpine3.10 AS build
RUN apk update && apk add git

# Copy source
WORKDIR /app/text-service
COPY . .

# Download dependencies application
WORKDIR /app/text-service/cmd
RUN go mod download

# Build application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.10 AS run

WORKDIR /etc/text-service/migrations
COPY ./resources/db/ .

WORKDIR /opt/app
COPY --from=build /app/text-service/cmd/cmd text-service
ENV GIN_MODE release
CMD ["./text-service"]
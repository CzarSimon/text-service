FROM golang:1.12-alpine3.10 AS build

# Copy source
WORKDIR /app/chatbot
COPY . .

# Download dependencies application
WORKDIR /app/chatbot/cmd
RUN go mod download

# Build application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.10 AS run

WORKDIR /etc/chatbot/migrations
COPY ./resources/db/ .

WORKDIR /opt/app
COPY --from=build /app/chatbot/cmd/cmd chatbot
ENV GIN_MODE release
CMD ["./chatbot"]
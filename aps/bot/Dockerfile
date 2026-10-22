FROM golang:1.26.3-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bot ./cmd

FROM alpine:3.23
WORKDIR /app
COPY --from=build /bot /app/bot
USER 10001:10001
CMD ["/app/bot"]

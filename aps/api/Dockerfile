FROM golang:1.26.3-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /api ./cmd

FROM alpine:3.23
WORKDIR /app
COPY --from=build /api /app/api
USER 10001:10001
EXPOSE 8080
CMD ["/app/api"]

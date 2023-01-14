# syntax=docker/dockerfile:1

FROM golang:1.16 AS builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY src/server.go ./
RUN go build -o wikimedia_api_server .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /build ./
EXPOSE 8080

CMD ["./wikimedia_api_server"]  
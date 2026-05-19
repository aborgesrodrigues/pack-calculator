FROM golang:1.26-alpine AS builder
WORKDIR /build
RUN apk add gcc git libc-dev

# Caching go dependencies.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pack-calculator ./cmd

FROM alpine:latest
WORKDIR /usr/local/bin/
EXPOSE 8080

COPY --from=builder /build/pack-calculator .
ENTRYPOINT ["./pack-calculator"]

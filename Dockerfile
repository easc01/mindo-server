FROM golang:1.24 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mindo-server ./cmd/main.go

# Install certs for scratch image
FROM alpine AS certs
RUN apk add --no-cache ca-certificates

FROM scratch

WORKDIR /root/

# Copy binary
COPY --from=builder /app/mindo-server .

# Copy CA certificates
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["./mindo-server"]

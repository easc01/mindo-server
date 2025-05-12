FROM golang:1.24 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mindo-server ./cmd/main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/mindo-server .

EXPOSE 8080

CMD ["./mindo-server"]

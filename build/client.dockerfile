FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o client ./cmd/client/client.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/client .

ENTRYPOINT ["./client"]

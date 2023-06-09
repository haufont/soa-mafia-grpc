FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server ./cmd/server/server.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/server .

ENTRYPOINT ["./server"]

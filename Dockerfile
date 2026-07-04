FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o e-wallet-system .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/e-wallet-system .

EXPOSE 8080

CMD ["./e-wallet-system"]

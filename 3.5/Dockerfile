FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Собираем в текущей директории
RUN cd cmd && go build -o ../server .

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/.env .env

EXPOSE 8080
CMD ["./server"]
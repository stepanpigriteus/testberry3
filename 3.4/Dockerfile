

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env
COPY assets /assets 
WORKDIR /app/cmd
RUN go build -o /server .

FROM debian:bookworm-slim

WORKDIR /app


COPY --from=builder /server .
COPY --from=builder /app/.env .env
COPY --from=builder /assets /assets  

EXPOSE 8080

CMD ["./server"]
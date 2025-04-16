FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./main.go 


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 1323

CMD ["./main"]
FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# Add a healthcheck for container monitoring
HEALTHCHECK --interval=30s --timeout=3s CMD wget -q -O- http://localhost:1323/health || exit 1
EXPOSE 1323
CMD ["./main"]
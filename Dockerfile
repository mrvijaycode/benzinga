
FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/.env ./

# Install ca-certificates to update the certificates store
RUN apt-get update && apt-get install -y ca-certificates

# Copy the custom SSL certificate to container. It's webhook.site ca signed certificate
COPY webhook.site.pem /usr/local/share/ca-certificates/webhook.site.pem

RUN update-ca-certificates

EXPOSE 8080

CMD ["./main"]

# Build stage
FROM golang:1.23-alpine AS builder

# Gerekli build araçlarını yükle
RUN apk add --no-cache git make

# Çalışma dizinini ayarla
WORKDIR /app

# Go modüllerini kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Binary'yi oluştur
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/payment-service ./cmd/payment

# Final stage
FROM alpine:3.19

# Gerekli paketleri yükle
RUN apk add --no-cache ca-certificates tzdata

# Çalışma dizinini ayarla
WORKDIR /app

# Binary'yi kopyala
COPY --from=builder /app/bin/payment-service .

# Environment variables
ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=paymentdb \
    HTTP_PORT=8080 \
    GRPC_PORT=50053 \
    KAFKA_BROKERS=kafka:9092

# Health check - HTTP ve gRPC için
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Port'ları aç
EXPOSE 8080 50053

# Graceful shutdown için sinyal yakalama
STOPSIGNAL SIGTERM

# Servisi başlat
CMD ["./payment-service"] 
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
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/customer-service ./cmd/customer

# Final stage
FROM alpine:3.19

# Gerekli paketleri yükle
RUN apk add --no-cache ca-certificates tzdata

# Çalışma dizinini ayarla
WORKDIR /app

# Binary'yi kopyala
COPY --from=builder /app/bin/customer-service .

# Environment variables
ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=customerdb \
    GRPC_PORT=50052

# Health check - gRPC için
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:50052/grpc.health.v1.Health/Check || exit 1

# Port'u aç
EXPOSE 50052

# Servisi başlat
CMD ["./customer-service"] 
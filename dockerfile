# Stage 1: Build binary
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Kompilasi CGO_ENABLED=0 menghasilkan binary statis murni
# Ldflags -w -s membuang debug info untuk efisiensi ruang (ukuran file lebih kecil)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o botengine main.go

# Stage 2: Eksekusi murni
FROM scratch
# Sertifikat SSL wajib disalin agar program bisa melakukan request HTTPS/REST API
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/botengine /botengine

# Eksekusi binary
CMD ["/botengine"]
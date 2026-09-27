```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY main.go .
COPY go.mod .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o gateway main.go

# Runtime stage - scratch (최소 이미지)
FROM scratch

COPY --from=builder /app/gateway /gateway
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["/gateway"]
```

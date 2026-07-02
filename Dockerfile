FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o refda-api ./cmd/api

FROM alpine:3.20
WORKDIR /app

# Copy CA bundle from builder (avoids apk/network in runtime image)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/refda-api .
COPY config.yaml .

RUN mkdir -p uploads
EXPOSE 8080
CMD ["./refda-api"]

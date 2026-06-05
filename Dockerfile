# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies separately
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]

FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 🔥 SEM GOARCH
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# ---------- runtime ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .
RUN chmod +x api

COPY .env .

EXPOSE 8000

CMD ["./api"]

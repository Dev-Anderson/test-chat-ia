FROM golang:1.23-alpine AS builder

WORKDIR /src

# Copia go.mod primeiro para cache de dependências
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia o restante do projeto
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Runtime
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /src/api /app/api
COPY .env /app/.env

EXPOSE 8000

CMD ["/app/api"]

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/url-shortener

# Runtime stage
FROM alpine:3.21

WORKDIR /app

COPY ui ./
COPY --from=builder /app/server .

EXPOSE 6969

CMD ["./server"]

# --- Builder ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o fileservice ./cmd/main.go

# --- Final image ---
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fileservice .
COPY migrations ./migrations
COPY config.yaml .
EXPOSE 8080
CMD ["./fileservice"]
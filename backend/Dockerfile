# Stage 1: builder
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Stage 2: runtime
FROM alpine:3.21
EXPOSE 8080
COPY --from=builder /app/server /app/server
CMD ["/app/server"]

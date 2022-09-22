# Build stage

FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8000
CMD ["/app/main"]
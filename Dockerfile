# Build Stage
FROM golang:1.21.0-bookworm AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .

# Generate Swagger docs
RUN swag i -g cmd/server/main.go -o docs

# Build Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server cmd/server/main.go

# Final Stage
FROM golang:1.21.0-bookworm

WORKDIR /app
COPY --from=builder /app/bin/server ./bin/server

EXPOSE 8001

CMD ["./bin/server"]

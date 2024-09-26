# Build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM gcr.io/distroless/static-debian12:latest-amd64

# Use dokku env for secrets
COPY --from=builder /app/.env /
# Dokku health check config
COPY --from=builder /app/app.json /

# Copy the binary from the builder stage
COPY --from=builder /app/main /



# Expose the port the app runs on
EXPOSE 80

# Run the binary
CMD ["/main"]
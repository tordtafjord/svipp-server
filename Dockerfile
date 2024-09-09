# Start from the official Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first and download dependencies
# This is done before copying the entire codebase to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o main ./cmd/api

# Expose the port the app runs on
EXPOSE 80

# Run the binary
CMD ["./main"]
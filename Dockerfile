# Use the official Golang image as the base image
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux cd ./cmd/api/ && go build -o f4rmhouse_api 

# Use a slim alpine image for the final container
FROM alpine:latest

# Install certificates (important for HTTPS connections)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/cmd/api/f4rmhouse_api .

# Expose the port the app runs on
EXPOSE 443 

# Command to run the executable
CMD ["./f4rmhouse_api"]

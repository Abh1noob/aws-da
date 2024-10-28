# Use the official Go Alpine image as the base image
FROM golang:alpine

# Update the package index, upgrade existing packages, and install Git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Install the Air live-reloading tool
RUN go install github.com/cosmtrek/air@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application code into the container
COPY . .

# Expose port 8080 to the host
EXPOSE 8080

# Command to run the Air live-reloading tool with the specified configuration
CMD ["air", "-c", ".air.toml"]

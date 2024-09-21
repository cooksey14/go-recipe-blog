# # Use the official Golang image to create a build environment
# FROM golang:1.18-alpine AS builder

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the source from the current directory to the Working Directory inside the container
# COPY . .

# # Build the Go app
# RUN go build -o main .

# # Start a new stage from scratch
# FROM alpine:latest

# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/main .

# # Copy config.json
# COPY config.json .

# # Command to run the executable
# CMD ["./main"]


# server/Dockerfile
FROM golang:1.20 AS builder

WORKDIR /app
COPY go.mod .          
COPY go.sum .          
RUN go mod tidy        
COPY . .            

# RUN GOOS=linux GOARCH=amd64 go build -o server .
RUN go build -o server .  
RUN ls -l /app



CMD ["./server"]
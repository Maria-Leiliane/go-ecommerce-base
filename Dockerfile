# Step Build
# Official Go image as a base for compiling the application.
FROM golang:1.25-alpine AS builder

# Defining the working directory within the container.
WORKDIR /app

# Copy module management files. go.sum may not exist
COPY go.mod go.sum* ./

# Download dependencies.
RUN go mod download

COPY . .

# Compiles the application, creating a static executable optimized for Linux.
# The -o flag defines the name of the output file.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./main.go

# Final Step
# Minimal base image (alpine), as Go tools are no longer needed.
FROM alpine:latest

# Defining the working directory within the container.
WORKDIR /app

# Copy ONLY the compiled executable from the build stage.
COPY --from=builder /app/main .

# Exposes port 8080 so we can connect to the API from outside the container.
EXPOSE 8080

# Command executed when the container starts.
CMD ["/app/main"]
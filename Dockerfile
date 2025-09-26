# --- Step 1: Compilation (Builder) ---
# Uses the official Go image with Alpine for a smaller build.
FROM golang:1.25-alpine AS builder

# Defining the working directory within the container.
WORKDIR /app

# Copies the module files first to take advantage of Docker's layer cache.
# 'go mod download' will only be rerun if go.mod/go.sum changes.
COPY go.mod go.sum* ./

# Download dependencies.
RUN go mod download

# Copies the rest of the source code.
COPY . .

# Statically compiles the application into a single executable.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./main.go

# --- Step 2: Final Image (Production) ---
# Uses a minimal image, as we only need the binary to run. FROM alpine:latest
FROM alpine:latest

# Defining the working directory within the container.
WORKDIR /app

# Copies only the compiled executable from the build step.
COPY --from=builder /app/main .

# Exposes the port the API will use.
EXPOSE 8080

# Command to start the application.
CMD ["/app/main"]
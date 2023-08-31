# Stage 1 - Build the base
FROM golang:1.21.0-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project including nested Go files
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/http

# Stage 2: Create a lightweight final image
FROM alpine:3.14.2

# Set the working directory
WORKDIR /app

# Set the environment variable
ENV GO_ENV=prod

# Copy the binary built in the previous stage
# COPY .env .
COPY --from=builder /main .

# Run
CMD ["./main"]

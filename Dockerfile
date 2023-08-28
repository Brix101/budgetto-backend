# Stage 1 - Build the base
FROM golang:1.21.0 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .


# Build
# RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/http/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/http

# Stage 2: Create a lightweight final image
FROM alpine:3.14.2

# Set the working directory
WORKDIR /app

# Copy the binary built in the previous stage
COPY .env .
COPY --from=builder /app/main .

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
# EXPOSE 5000

# Run
CMD ["./main"]

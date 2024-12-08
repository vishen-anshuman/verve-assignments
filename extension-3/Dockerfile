# syntax=docker/dockerfile:1

# Step 1: Use a Golang image to build the application
FROM golang:1.22 AS build

WORKDIR /app

# Install Git to fetch the code from GitHub
RUN apt-get update && apt-get install -y git

# Clone the repository from GitHub
RUN git clone https://github.com/vishen-anshuman/verve-assignments.git .

WORKDIR /app/extension-3

# Download Go dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 go build -o /bin/server ./server.go

# Step 2: Create the runtime image
FROM alpine:3.17.2 AS final

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates

# Create a non-privileged user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# Set up the logs directory with correct permissions
RUN mkdir -p /app/logs && chown appuser:appuser /app/logs

USER appuser

# Copy the built Go executable from the build stage
COPY --from=build /bin/server /bin/

# Expose the application port
EXPOSE 8080

# Set the entrypoint for the container
ENTRYPOINT ["/bin/server"]


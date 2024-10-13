# Use a more suitable base image for Alpine
FROM golang:1.23-alpine

# Install necessary packages including bash, curl, and dcron (Alpine's cron package)
RUN apk add --no-cache bash curl dcron

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o /app/scraper ./cmd/scraper/main.go

# Copy the cron file to the container
COPY cronfile /etc/cron.d/cronfile

# Give execution rights on the cron job
RUN chmod 0644 /etc/cron.d/cronfile

# Register the cron job
RUN crontab /etc/cron.d/cronfile

# Copy entrypoint.sh script
COPY entrypoint.sh /app/entrypoint.sh

# Make sure the script is executable
RUN chmod +x /app/entrypoint.sh

# Use the entrypoint.sh script as the entry point
ENTRYPOINT ["/app/entrypoint.sh"]

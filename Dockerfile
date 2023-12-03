# Use the Golang base image version 1.21
FROM golang:1.21

# Set the working directory
WORKDIR /app

# Install dependencies for downloading and installing the migrate utility
RUN apt-get update && apt-get install -y wget && \
    wget -O /usr/local/bin/migrate https://github.com/golang-migrate/migrate/releases/download/v4.16.0/migrate.linux-amd64.tar.gz && \
    tar -xzf /usr/local/bin/migrate -C /usr/local/bin/ && \
    chmod +x /usr/local/bin/migrate && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy go.mod and go.sum files and download dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# Copy all other project files into the current directory
COPY . .

# Build the main binary file
RUN go build -o main ./cmd/api/

# Apply migrations and run the application on container start
CMD sh -c "migrate -path ./migrations -database postgres://gifts:gifts@db:5432/gifts?sslmode=disable up && ./main"

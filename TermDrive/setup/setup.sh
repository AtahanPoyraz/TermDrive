#!/bin/bash

# Starts Docker Compose build and run commands
echo "Starting Docker Compose setup..."

# Stop and remove all containers, networks, and volumes defined in the docker-compose file
docker compose down -v

# Build and run services
docker compose up --build -d

# Check the status
echo "Services are up and running:"
docker compose ps

# Updating packages
echo "Running go mod tidy..."
go mod tidy

echo "Go modules updated."
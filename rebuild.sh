#!/bin/bash

# To run this script in WSL2, make sure you have Docker Desktop configured to use the WSL2 backend.
# I have it configured to use Ubuntu, when I run docker compose up, wsl runs the windows command 
# and exposes the windows ports.
# Docker Desktop settings -> Resources -> WSL Integration -> Enable integration with my default WSL distro

# Exit immediately if a command exits with a non-zero status
set -e 

VOLUME_NAME="chronocode_postgres_data"

# Stop and remove existing containers
echo "Stopping and removing existing containers..."
docker compose down

# Remove existing database volume
echo "Removing existing database volume..."
docker volume rm -f $VOLUME_NAME || true

# Rebuild and start containers
echo "Rebuilding and starting containers..."
docker compose build

# Start containers in detached mode
echo "Starting containers..."
docker compose up -d

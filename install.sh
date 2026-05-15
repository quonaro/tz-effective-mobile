#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_IMAGE="tz-backend"
FRONTEND_IMAGE="tz-frontend"

echo "Loading Docker images..."

if [ -f "$SCRIPT_DIR/${BACKEND_IMAGE}.tar.gz" ]; then
    echo "Loading backend image..."
    gunzip -c "$SCRIPT_DIR/${BACKEND_IMAGE}.tar.gz" | docker load
else
    echo "Error: ${BACKEND_IMAGE}.tar.gz not found"
    exit 1
fi

if [ -f "$SCRIPT_DIR/${FRONTEND_IMAGE}.tar.gz" ]; then
    echo "Loading frontend image..."
    gunzip -c "$SCRIPT_DIR/${FRONTEND_IMAGE}.tar.gz" | docker load
else
    echo "Error: ${FRONTEND_IMAGE}.tar.gz not found"
    exit 1
fi

echo "Images loaded successfully!"
echo ""
echo "To run the application:"
echo "  docker compose up -d"
echo ""
echo "Or with custom .env file:"
echo "  docker compose --env-file .env up -d"

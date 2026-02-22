#!/bin/bash
# uduXPass Backend Startup Script
# This script ensures the backend starts with correct PostgreSQL connection

set -e

echo "🚀 Starting uduXPass Backend..."

# Set correct DATABASE_URL for PostgreSQL (Unix socket)
export DATABASE_URL="host=/var/run/postgresql user=postgres dbname=uduxpass sslmode=disable"

# Start backend in background
cd /home/ubuntu/uduxpass-backend
nohup ./uduxpass-api > /tmp/uduxpass-backend.log 2>&1 &

# Wait for backend to start
echo "⏳ Waiting for backend to start..."
for i in {1..15}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Backend started successfully!"
        curl -s http://localhost:8080/health | python3 -m json.tool
        exit 0
    fi
    sleep 1
done

echo "❌ Backend failed to start. Check logs:"
tail -20 /tmp/uduxpass-backend.log
exit 1

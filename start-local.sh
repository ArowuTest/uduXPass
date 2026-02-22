#!/bin/bash

# uduXPass Platform - Local Deployment Startup Script
# This script starts all services and initializes the database

set -e

echo "🚀 Starting uduXPass Platform (Local Deployment)"
echo "================================================"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if .env file exists, if not create from example
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please update .env with your SMTP and Paystack credentials"
    echo ""
fi

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
docker-compose down

# Build and start all services
echo ""
echo "🔨 Building Docker images..."
docker-compose build

echo ""
echo "🚀 Starting all services..."
docker-compose up -d

# Wait for database to be ready
echo ""
echo "⏳ Waiting for database to be ready..."
sleep 10

# Run migrations and seed data
echo ""
echo "🗄️  Running database migrations..."
docker-compose exec -T backend ./uduxpass-api migrate || echo "⚠️  Migrations may have already run"

echo ""
echo "🌱 Seeding database with initial data..."
docker-compose exec -T backend ./uduxpass-api seed || echo "⚠️  Seed data may already exist"

# Show status
echo ""
echo "✅ uduXPass Platform is now running!"
echo "================================================"
echo ""
echo "📱 Access the applications:"
echo "   - Customer Frontend: http://localhost:3000"
echo "   - Scanner App:       http://localhost:3001"
echo "   - Backend API:       http://localhost:8080"
echo "   - Database:          localhost:5432"
echo ""
echo "👤 Test Credentials:"
echo "   Admin:    admin@uduxpass.com / Admin123!"
echo "   Scanner:  scanner@uduxpass.com / Scanner123!"
echo "   Customer: customer@uduxpass.com / Customer123!"
echo ""
echo "📊 View logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Stop all services:"
echo "   docker-compose down"
echo ""
echo "================================================"

#!/bin/bash

# uduXPass Platform - Stop Script

echo "🛑 Stopping uduXPass Platform..."
docker-compose down

echo ""
echo "✅ All services stopped!"
echo ""
echo "💡 To remove all data (including database):"
echo "   docker-compose down -v"

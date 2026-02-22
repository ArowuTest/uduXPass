#!/bin/bash
cd /home/ubuntu/uduxpass-backend
export DATABASE_URL="host=/var/run/postgresql user=postgres dbname=uduxpass sslmode=disable"
export CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:3003,http://localhost:5173,http://localhost:5174"
exec ./uduxpass-api

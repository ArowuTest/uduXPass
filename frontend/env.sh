#!/bin/sh

# Replace environment variables in built files
# This allows runtime configuration of the frontend

# Default values
API_URL=${API_URL:-"http://localhost:8080"}
APP_ENV=${APP_ENV:-"production"}

# Replace placeholders in JavaScript files
find /usr/share/nginx/html -name "*.js" -exec sed -i "s|__API_URL__|$API_URL|g" {} \;
find /usr/share/nginx/html -name "*.js" -exec sed -i "s|__APP_ENV__|$APP_ENV|g" {} \;

echo "Environment variables injected:"
echo "API_URL: $API_URL"
echo "APP_ENV: $APP_ENV"


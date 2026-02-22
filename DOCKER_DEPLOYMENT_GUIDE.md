# uduXPass Platform - Docker Local Deployment Guide

**Version:** 1.0  
**Date:** February 20, 2026  
**Status:** Production Ready

---

## 🚀 Quick Start (One Command)

```bash
./start-local.sh
```

That's it! The entire platform will be running in Docker containers.

---

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [Architecture](#architecture)
4. [Services](#services)
5. [Configuration](#configuration)
6. [Usage](#usage)
7. [Troubleshooting](#troubleshooting)
8. [Advanced](#advanced)

---

## Prerequisites

### Required Software:
- **Docker** 20.10+ ([Install Docker](https://docs.docker.com/get-docker/))
- **Docker Compose** 2.0+ (included with Docker Desktop)

### System Requirements:
- **RAM:** 4GB minimum, 8GB recommended
- **Disk Space:** 10GB free space
- **OS:** Linux, macOS, or Windows with WSL2

### Verify Installation:
```bash
docker --version
docker-compose --version
```

---

## Quick Start

### 1. Extract Repository
```bash
tar -xzf uduxpass-fullstack-fixed-final-20260220.tar.gz
cd uduxpass-fullstack-fixed-final-20260220
```

### 2. Configure Environment (Optional)
```bash
# Copy example environment file
cp .env.example .env

# Edit with your credentials (optional for local testing)
nano .env
```

### 3. Start Platform
```bash
./start-local.sh
```

### 4. Access Applications

**Customer Frontend:** http://localhost:3000  
**Scanner App:** http://localhost:3001  
**Backend API:** http://localhost:8080  
**Database:** localhost:5432

### 5. Test Credentials

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@uduxpass.com | Admin123! |
| Scanner | scanner@uduxpass.com | Scanner123! |
| Customer | customer@uduxpass.com | Customer123! |

---

## Architecture

### Docker Services

```
┌─────────────────────────────────────────────────┐
│                                                 │
│  Customer Frontend (Port 3000)                  │
│  ┌────────────────────────────────────────┐    │
│  │  React + Vite + Nginx                  │    │
│  └────────────────────────────────────────┘    │
│                                                 │
│  Scanner App (Port 3001)                        │
│  ┌────────────────────────────────────────┐    │
│  │  React + Vite + Nginx                  │    │
│  └────────────────────────────────────────┘    │
│                                                 │
│  Backend API (Port 8080)                        │
│  ┌────────────────────────────────────────┐    │
│  │  Go + Gin + JWT                        │    │
│  └────────────────────────────────────────┘    │
│                                                 │
│  PostgreSQL Database (Port 5432)                │
│  ┌────────────────────────────────────────┐    │
│  │  PostgreSQL 14                         │    │
│  └────────────────────────────────────────┘    │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Network
All services communicate through a private Docker network: `uduxpass-network`

### Volumes
- `postgres_data` - Persistent database storage

---

## Services

### 1. PostgreSQL Database

**Container:** `uduxpass-db`  
**Image:** `postgres:14-alpine`  
**Port:** 5432  
**Credentials:**
- User: `uduxpass`
- Password: `uduxpass_local_password`
- Database: `uduxpass_db`

**Health Check:** Runs every 10 seconds

### 2. Backend API

**Container:** `uduxpass-backend`  
**Build:** `./backend/Dockerfile`  
**Port:** 8080  
**Features:**
- User authentication (JWT)
- Event management
- Ticket management
- Payment processing (Paystack)
- Email notifications (SMTP)

**Health Check:** `http://localhost:8080/health`

### 3. Customer Frontend

**Container:** `uduxpass-frontend`  
**Build:** `./frontend/Dockerfile`  
**Port:** 3000  
**Tech Stack:**
- React 19
- TypeScript
- Vite
- Nginx

**Health Check:** `http://localhost:3000/`

### 4. Scanner App

**Container:** `uduxpass-scanner`  
**Build:** `./uduxpass-scanner-app/Dockerfile`  
**Port:** 3001  
**Tech Stack:**
- React 19
- TypeScript
- Vite
- Nginx

**Health Check:** `http://localhost:3001/`

---

## Configuration

### Environment Variables

#### Backend (.env or docker-compose.yml)
```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=uduxpass
DB_PASSWORD=uduxpass_local_password
DB_NAME=uduxpass_db

# JWT
JWT_SECRET=local_development_jwt_secret_key

# Email (Optional)
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Payment (Test Keys)
PAYSTACK_SECRET_KEY=sk_test_your_test_key
PAYSTACK_PUBLIC_KEY=pk_test_your_test_key
```

#### Frontend
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_PAYSTACK_PUBLIC_KEY=pk_test_your_test_key
```

#### Scanner App
```env
VITE_API_BASE_URL=http://localhost:8080
```

---

## Usage

### Start All Services
```bash
./start-local.sh
```

### Stop All Services
```bash
./stop-local.sh
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f scanner
docker-compose logs -f postgres
```

### Check Service Status
```bash
docker-compose ps
```

### Restart a Service
```bash
docker-compose restart backend
docker-compose restart frontend
docker-compose restart scanner
```

### Access Database
```bash
# Using psql
docker-compose exec postgres psql -U uduxpass -d uduxpass_db

# Using external tool
Host: localhost
Port: 5432
User: uduxpass
Password: uduxpass_local_password
Database: uduxpass_db
```

### Run Migrations
```bash
docker-compose exec backend ./uduxpass-api migrate
```

### Seed Database
```bash
docker-compose exec backend ./uduxpass-api seed
```

### Rebuild Services
```bash
# Rebuild all
docker-compose build

# Rebuild specific service
docker-compose build backend
docker-compose build frontend
docker-compose build scanner

# Rebuild and restart
docker-compose up -d --build
```

---

## Troubleshooting

### Port Already in Use

**Problem:** Port 3000, 3001, 8080, or 5432 is already in use

**Solution:**
```bash
# Find process using port
lsof -i :3000
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in docker-compose.yml
ports:
  - "3002:80"  # Use different host port
```

### Container Won't Start

**Problem:** Container exits immediately

**Solution:**
```bash
# Check logs
docker-compose logs backend

# Check container status
docker-compose ps

# Rebuild container
docker-compose build backend
docker-compose up -d backend
```

### Database Connection Failed

**Problem:** Backend can't connect to database

**Solution:**
```bash
# Check if database is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Restart database
docker-compose restart postgres

# Wait for database to be ready
docker-compose up -d postgres
sleep 10
docker-compose up -d backend
```

### Frontend Can't Reach Backend

**Problem:** API calls fail with CORS or network errors

**Solution:**
1. Check backend is running: `docker-compose ps backend`
2. Check backend logs: `docker-compose logs backend`
3. Verify CORS settings in backend
4. Check VITE_API_BASE_URL in frontend

### Build Errors

**Problem:** Docker build fails

**Solution:**
```bash
# Clean Docker cache
docker system prune -a

# Rebuild from scratch
docker-compose build --no-cache

# Check Dockerfile syntax
docker-compose config
```

### Out of Disk Space

**Problem:** Docker runs out of space

**Solution:**
```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove everything
docker system prune -a --volumes
```

---

## Advanced

### Production Deployment

For production deployment, see `DEPLOYMENT_GUIDE.md` for:
- SSL/HTTPS configuration
- Nginx reverse proxy
- Environment-specific settings
- Security hardening
- Monitoring and logging

### Custom Configuration

#### Change Database Password
1. Update `docker-compose.yml`:
```yaml
environment:
  POSTGRES_PASSWORD: your_new_password
  DB_PASSWORD: your_new_password
```

2. Remove existing volume:
```bash
docker-compose down -v
docker-compose up -d
```

#### Add Custom Domain
1. Update `/etc/hosts`:
```
127.0.0.1 uduxpass.local
```

2. Access at: http://uduxpass.local:3000

#### Enable HTTPS (Local)
1. Generate self-signed certificate
2. Update Nginx configuration
3. Update docker-compose.yml ports

### Development Mode

For active development with hot reload:

```bash
# Start only database
docker-compose up -d postgres

# Run backend locally
cd backend
go run cmd/main.go

# Run frontend locally
cd frontend
pnpm run dev

# Run scanner locally
cd uduxpass-scanner-app/client
pnpm run dev
```

### Backup and Restore

#### Backup Database
```bash
docker-compose exec postgres pg_dump -U uduxpass uduxpass_db > backup.sql
```

#### Restore Database
```bash
cat backup.sql | docker-compose exec -T postgres psql -U uduxpass -d uduxpass_db
```

### Performance Tuning

#### Increase PostgreSQL Memory
```yaml
postgres:
  command: postgres -c shared_buffers=256MB -c max_connections=200
```

#### Enable Nginx Caching
Update `frontend/nginx.conf` and `uduxpass-scanner-app/nginx.conf`

---

## Docker Commands Reference

### Basic Commands
```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f

# Check status
docker-compose ps

# Restart service
docker-compose restart <service>

# Rebuild service
docker-compose build <service>

# Execute command in container
docker-compose exec <service> <command>
```

### Cleanup Commands
```bash
# Remove stopped containers
docker-compose down

# Remove volumes (deletes data!)
docker-compose down -v

# Remove images
docker-compose down --rmi all

# Clean everything
docker system prune -a --volumes
```

---

## Support

### Documentation
- **E2E Tests:** `E2E_TEST_REPORT.md`
- **Production Deployment:** `DEPLOYMENT_GUIDE.md`
- **Validation Fixes:** `VALIDATION_FIXES_SUMMARY.md`
- **Package Contents:** `PACKAGE_CONTENTS.md`

### Common Issues
- Check logs: `docker-compose logs -f`
- Verify services: `docker-compose ps`
- Check health: `docker-compose ps` (look for "healthy" status)

### Getting Help
1. Check logs for error messages
2. Verify all prerequisites are installed
3. Ensure ports are not in use
4. Try rebuilding: `docker-compose build --no-cache`

---

## Conclusion

The uduXPass platform is fully dockerized and ready for local deployment. With a single command (`./start-local.sh`), you can have the entire platform running locally for development and testing.

**Happy Developing! 🚀**

---

**Prepared by:** Manus AI Agent  
**Date:** February 20, 2026  
**Version:** Docker Local Deployment v1.0

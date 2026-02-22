# uduXPass Platform - Deployment Guide

**Version:** 1.0 (Fixed Release)  
**Date:** February 20, 2026  
**Status:** Production Ready

---

## Table of Contents

1. [Overview](#overview)
2. [Repository Structure](#repository-structure)
3. [Prerequisites](#prerequisites)
4. [Backend Deployment](#backend-deployment)
5. [Frontend Deployment](#frontend-deployment)
6. [Scanner App Deployment](#scanner-app-deployment)
7. [Environment Configuration](#environment-configuration)
8. [Database Setup](#database-setup)
9. [Testing](#testing)
10. [Troubleshooting](#troubleshooting)

---

## Overview

The uduXPass platform consists of three main applications:

1. **Backend API** (Go) - Port 8080
2. **Customer Frontend** (React + TypeScript) - Port 5173
3. **Scanner App** (React + TypeScript) - Port 3000

All applications are production-ready and fully tested.

---

## Repository Structure

```
uduxpass-fullstack-fixed-20260220/
├── backend/                    # Go backend API
│   ├── cmd/
│   │   └── main.go            # Entry point
│   ├── internal/
│   │   ├── api/               # API handlers
│   │   ├── models/            # Database models
│   │   ├── services/          # Business logic
│   │   └── middleware/        # Middleware
│   ├── config/
│   │   └── config.go          # Configuration
│   ├── migrations/            # Database migrations
│   ├── seeds/                 # Seed data
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # Customer-facing React app
│   ├── src/
│   │   ├── pages/             # Page components
│   │   ├── components/        # Reusable components
│   │   ├── contexts/          # React contexts
│   │   ├── services/          # API services
│   │   └── App.tsx
│   ├── package.json
│   └── vite.config.ts
│
├── uduxpass-scanner-app/      # Scanner PWA
│   ├── client/
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   ├── components/
│   │   │   └── App.tsx
│   │   └── package.json
│   └── package.json
│
├── E2E_TEST_REPORT.md         # E2E test results
├── UDUXPASS_FINAL_PRODUCTION_REPORT.md  # Production report
├── VALIDATION_FIXES_SUMMARY.md          # Fix summary
└── DEPLOYMENT_GUIDE.md        # This file
```

---

## Prerequisites

### Required Software:
- **Go** 1.21 or higher
- **Node.js** 18.x or higher
- **pnpm** 8.x or higher
- **PostgreSQL** 14.x or higher
- **Git**

### Optional:
- **Docker** (for containerized deployment)
- **Nginx** (for reverse proxy)
- **SSL Certificate** (for HTTPS)

---

## Backend Deployment

### 1. Extract Repository
```bash
tar -xzf uduxpass-fullstack-fixed-20260220.tar.gz
cd backend
```

### 2. Install Dependencies
```bash
go mod download
go mod verify
```

### 3. Configure Environment
```bash
cp .env.example .env
nano .env
```

**Required Environment Variables:**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=uduxpass
DB_PASSWORD=your_secure_password
DB_NAME=uduxpass_db
DB_SSLMODE=disable  # Use 'require' in production

# Server
PORT=8080
ENV=production
FRONTEND_URL=https://yourdomain.com

# JWT
JWT_SECRET=your_jwt_secret_key_minimum_32_characters_long

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@uduxpass.com

# Payment (Paystack)
PAYSTACK_SECRET_KEY=sk_live_your_paystack_secret_key
PAYSTACK_PUBLIC_KEY=pk_live_your_paystack_public_key

# MoMo PSB (Optional)
MOMO_API_KEY=your_momo_api_key
MOMO_API_SECRET=your_momo_api_secret
MOMO_BASE_URL=https://api.momo.com

# CORS
ALLOWED_ORIGINS=https://yourdomain.com,https://scanner.yourdomain.com
```

### 4. Setup Database
```bash
# Create database
createdb uduxpass_db

# Run migrations
go run cmd/main.go migrate

# Seed data
go run cmd/main.go seed
```

### 5. Build and Run
```bash
# Build
go build -o uduxpass-api cmd/main.go

# Run
./uduxpass-api
```

### 6. Verify Backend
```bash
curl http://localhost:8080/health
# Expected: {"status":"ok","timestamp":"..."}
```

---

## Frontend Deployment

### 1. Navigate to Frontend
```bash
cd ../frontend
```

### 2. Install Dependencies
```bash
pnpm install
```

### 3. Configure Environment
```bash
cp .env.example .env
nano .env
```

**Required Environment Variables:**
```env
VITE_API_BASE_URL=https://api.yourdomain.com
VITE_PAYSTACK_PUBLIC_KEY=pk_live_your_paystack_public_key
```

### 4. Build for Production
```bash
pnpm run build
```

### 5. Deploy Build
```bash
# Option 1: Serve with nginx
sudo cp -r dist/* /var/www/html/

# Option 2: Use a static hosting service
# Upload dist/ folder to Vercel, Netlify, etc.

# Option 3: Serve with Node.js
pnpm add -g serve
serve -s dist -p 5173
```

### 6. Verify Frontend
```bash
curl http://localhost:5173
# Expected: HTML content
```

---

## Scanner App Deployment

### 1. Navigate to Scanner App
```bash
cd ../uduxpass-scanner-app
```

### 2. Install Dependencies
```bash
pnpm install
```

### 3. Configure Environment
```bash
cp .env.example .env
nano .env
```

**Required Environment Variables:**
```env
VITE_API_BASE_URL=https://api.yourdomain.com
```

### 4. Build for Production
```bash
pnpm run build
```

### 5. Deploy Build
```bash
# Serve on port 3000
pnpm add -g serve
serve -s dist -p 3000
```

### 6. Verify Scanner App
```bash
curl http://localhost:3000
# Expected: HTML content
```

---

## Environment Configuration

### Production Environment Variables

#### Backend (.env)
```env
# Database
DB_HOST=your-db-host.com
DB_PORT=5432
DB_USER=uduxpass_prod
DB_PASSWORD=strong_production_password
DB_NAME=uduxpass_production
DB_SSLMODE=require

# Server
PORT=8080
ENV=production
FRONTEND_URL=https://uduxpass.com

# JWT (Generate with: openssl rand -base64 32)
JWT_SECRET=your_production_jwt_secret_minimum_32_characters

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@uduxpass.com
SMTP_PASSWORD=your_app_specific_password
SMTP_FROM=noreply@uduxpass.com

# Paystack
PAYSTACK_SECRET_KEY=sk_live_xxxxxxxxxxxxx
PAYSTACK_PUBLIC_KEY=pk_live_xxxxxxxxxxxxx

# CORS
ALLOWED_ORIGINS=https://uduxpass.com,https://scanner.uduxpass.com
```

#### Frontend (.env)
```env
VITE_API_BASE_URL=https://api.uduxpass.com
VITE_PAYSTACK_PUBLIC_KEY=pk_live_xxxxxxxxxxxxx
```

#### Scanner App (.env)
```env
VITE_API_BASE_URL=https://api.uduxpass.com
```

---

## Database Setup

### 1. Create Production Database
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and user
CREATE DATABASE uduxpass_production;
CREATE USER uduxpass_prod WITH ENCRYPTED PASSWORD 'strong_password';
GRANT ALL PRIVILEGES ON DATABASE uduxpass_production TO uduxpass_prod;
\q
```

### 2. Run Migrations
```bash
cd backend
go run cmd/main.go migrate
```

### 3. Seed Initial Data
```bash
go run cmd/main.go seed
```

**Seeded Data Includes:**
- 4 Events (Burna Boy, Wizkid, Davido, Afro Nation)
- 3 Ticket tiers per event
- 1 Admin user (admin@uduxpass.com / Admin123!)
- 1 Scanner user (scanner@uduxpass.com / Scanner123!)
- 1 Test customer (customer@uduxpass.com / Customer123!)

### 4. Verify Database
```bash
psql -U uduxpass_prod -d uduxpass_production -c "SELECT COUNT(*) FROM events;"
# Expected: 4
```

---

## Testing

### 1. Backend API Tests
```bash
cd backend

# Test health endpoint
curl http://localhost:8080/health

# Test registration
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "User",
    "email": "test@example.com",
    "phone": "+2348099999999",
    "password": "TestPassword123!"
  }'

# Test login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@uduxpass.com",
    "password": "Admin123!"
  }'

# Test events list
curl http://localhost:8080/v1/events
```

### 2. Frontend Tests
```bash
# Open in browser
open http://localhost:5173

# Test pages:
# - Home page: http://localhost:5173/
# - Events: http://localhost:5173/events
# - Event details: http://localhost:5173/events/1
# - Registration: http://localhost:5173/register
# - Login: http://localhost:5173/login
```

### 3. Scanner App Tests
```bash
# Open in browser
open http://localhost:3000

# Test login with scanner credentials:
# Email: scanner@uduxpass.com
# Password: Scanner123!
```

### 4. E2E Flow Test
1. Register new user
2. Browse events
3. View event details
4. Add tickets to cart
5. Checkout (test mode)
6. Receive email confirmation
7. Login to scanner app
8. Scan ticket QR code

---

## Troubleshooting

### Backend Issues

#### Database Connection Failed
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check connection
psql -U uduxpass_prod -d uduxpass_production -c "SELECT 1;"

# Check .env file
cat .env | grep DB_
```

#### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

#### Migration Errors
```bash
# Reset database (CAUTION: Deletes all data)
dropdb uduxpass_production
createdb uduxpass_production
go run cmd/main.go migrate
go run cmd/main.go seed
```

### Frontend Issues

#### API Connection Failed
```bash
# Check VITE_API_BASE_URL in .env
cat .env | grep VITE_API_BASE_URL

# Test API directly
curl https://api.yourdomain.com/health

# Check browser console for CORS errors
```

#### Build Errors
```bash
# Clear cache and reinstall
rm -rf node_modules pnpm-lock.yaml
pnpm install
pnpm run build
```

#### Phone Validation Errors
```bash
# Verify fixes are applied
grep -A 5 "phoneRegex" src/pages/auth/RegisterPage.tsx

# Should show: /^(\+234|234|0)\d{10}$/
```

### Scanner App Issues

#### Login Failed
```bash
# Verify scanner user exists in database
psql -U uduxpass_prod -d uduxpass_production \
  -c "SELECT email, role FROM users WHERE role='scanner';"

# Reset scanner password
psql -U uduxpass_prod -d uduxpass_production \
  -c "UPDATE users SET password='$2a$10$...' WHERE email='scanner@uduxpass.com';"
```

---

## Production Deployment Checklist

### Pre-Deployment:
- [ ] All tests passing
- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] Seed data loaded
- [ ] SSL certificates installed
- [ ] Domain DNS configured
- [ ] SMTP credentials verified
- [ ] Paystack keys verified

### Deployment:
- [ ] Backend deployed and running
- [ ] Frontend built and deployed
- [ ] Scanner app built and deployed
- [ ] Nginx/reverse proxy configured
- [ ] HTTPS enabled
- [ ] CORS configured correctly

### Post-Deployment:
- [ ] Health checks passing
- [ ] Registration flow tested
- [ ] Login flow tested
- [ ] Event browsing tested
- [ ] Ticket purchase tested (test mode)
- [ ] Scanner app tested
- [ ] Email delivery tested
- [ ] Error logging configured
- [ ] Monitoring setup

---

## Performance Optimization

### Backend:
- Enable database connection pooling
- Add Redis for session caching
- Enable gzip compression
- Configure rate limiting

### Frontend:
- Enable CDN for static assets
- Configure browser caching
- Optimize images
- Enable lazy loading

### Database:
- Add indexes on frequently queried fields
- Enable query caching
- Configure connection pooling
- Regular VACUUM and ANALYZE

---

## Security Checklist

- [ ] All passwords hashed with bcrypt
- [ ] JWT tokens with expiration
- [ ] HTTPS enabled
- [ ] CORS properly configured
- [ ] SQL injection protection
- [ ] XSS protection
- [ ] CSRF protection
- [ ] Rate limiting enabled
- [ ] Input validation on all endpoints
- [ ] Secure headers configured

---

## Monitoring and Logging

### Recommended Tools:
- **Application Monitoring:** New Relic, Datadog
- **Error Tracking:** Sentry
- **Log Management:** ELK Stack, Papertrail
- **Uptime Monitoring:** Pingdom, UptimeRobot

### Key Metrics to Monitor:
- API response times
- Error rates
- Database query performance
- User registration rate
- Ticket sales
- Scanner app usage

---

## Support and Documentation

### Documentation:
- E2E Test Report: `E2E_TEST_REPORT.md`
- Production Report: `UDUXPASS_FINAL_PRODUCTION_REPORT.md`
- Validation Fixes: `VALIDATION_FIXES_SUMMARY.md`
- Deployment Guide: `DEPLOYMENT_GUIDE.md` (this file)

### Contact:
- Technical Support: support@uduxpass.com
- Emergency Hotline: +234 (0) 800 UDUXPASS

---

## Conclusion

The uduXPass platform is production-ready and fully tested. Follow this guide for a smooth deployment process.

**Good luck with your launch! 🚀**

---

**Prepared by:** Manus AI Agent  
**Date:** February 20, 2026  
**Version:** Production Release v1.0

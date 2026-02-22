# uduXPass Platform - Deployment Checklist

**Repository:** https://github.com/ArowuTest/uduXPass  
**Version:** 1.0.0  
**Date:** February 22, 2026

---

## Pre-Deployment Checklist

### 1. Environment Setup

#### Database
- [ ] PostgreSQL 14+ installed
- [ ] Database created: `uduxpass_db`
- [ ] Database user created with appropriate permissions
- [ ] Database connection string configured
- [ ] SSL/TLS enabled for database connections

#### Backend Server
- [ ] Go 1.21+ installed
- [ ] Environment variables configured (see `.env.example`)
- [ ] SMTP credentials configured (for email)
- [ ] Paystack API keys configured (for payments)
- [ ] MoMo API credentials configured (optional)
- [ ] JWT secret generated (secure random string)
- [ ] CORS origins configured for frontend URL

#### Frontend Server
- [ ] Node.js 20+ installed
- [ ] pnpm installed
- [ ] Environment variables configured
- [ ] API endpoint URL configured
- [ ] Build directory configured
- [ ] CDN/hosting configured (if applicable)

#### Scanner App
- [ ] Node.js 20+ installed
- [ ] pnpm installed
- [ ] Environment variables configured
- [ ] API endpoint URL configured
- [ ] PWA manifest configured
- [ ] Service worker configured

---

## Deployment Steps

### Step 1: Clone Repository

```bash
git clone https://github.com/ArowuTest/uduXPass.git
cd uduXPass
```

**Verification:**
- [ ] Repository cloned successfully
- [ ] All files present (backend, frontend, scanner, docs)
- [ ] `.gitignore` present

---

### Step 2: Database Setup

#### 2.1 Create Database

```bash
sudo -u postgres psql
CREATE DATABASE uduxpass_db;
CREATE USER uduxpass WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE uduxpass_db TO uduxpass;
\q
```

**Verification:**
- [ ] Database created
- [ ] User created with permissions
- [ ] Can connect to database

#### 2.2 Run Migrations

```bash
cd backend
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=uduxpass
export DB_PASSWORD=your_secure_password
export DB_NAME=uduxpass_db

# Run migrations
go run cmd/migrate/main.go
```

**Verification:**
- [ ] All 11 migrations applied successfully
- [ ] Tables created (users, events, tickets, orders, etc.)
- [ ] No migration errors

#### 2.3 Load Seed Data (Optional for Production)

```bash
# For development/testing only
psql -U uduxpass -d uduxpass_db -f migrations/009_comprehensive_seed_data.sql
```

**Verification:**
- [ ] Test users created (if loading seed data)
- [ ] Test events created (if loading seed data)
- [ ] Ticket tiers configured (if loading seed data)

---

### Step 3: Backend Deployment

#### 3.1 Configure Environment

Create `.env` file in `backend/` directory:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=uduxpass
DB_PASSWORD=your_secure_password
DB_NAME=uduxpass_db
DB_SSLMODE=require

# Server
PORT=8080
FRONTEND_URL=https://your-frontend-domain.com

# JWT
JWT_SECRET=your_secure_jwt_secret_minimum_32_characters

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@uduxpass.com

# Paystack
PAYSTACK_SECRET_KEY=sk_live_your_paystack_secret_key
PAYSTACK_PUBLIC_KEY=pk_live_your_paystack_public_key

# MoMo (Optional)
MOMO_API_USER=your_momo_api_user
MOMO_API_KEY=your_momo_api_key
MOMO_SUBSCRIPTION_KEY=your_momo_subscription_key
```

**Verification:**
- [ ] All required environment variables set
- [ ] Secure secrets generated
- [ ] API keys configured
- [ ] CORS origins correct

#### 3.2 Build Backend

```bash
cd backend
go build -o uduxpass-api cmd/api/main.go
```

**Verification:**
- [ ] Build successful
- [ ] Binary created: `uduxpass-api`
- [ ] No build errors

#### 3.3 Start Backend

```bash
# Production (with systemd)
sudo systemctl start uduxpass-api
sudo systemctl enable uduxpass-api

# Or manually
./uduxpass-api
```

**Verification:**
- [ ] Backend started successfully
- [ ] Listening on port 8080
- [ ] Health check endpoint responds: `curl http://localhost:8080/health`
- [ ] No startup errors in logs

---

### Step 4: Frontend Deployment

#### 4.1 Configure Environment

Create `.env` file in `frontend/` directory:

```env
VITE_API_URL=https://api.your-domain.com
VITE_APP_NAME=uduXPass
VITE_APP_DESCRIPTION=Premium Event Ticketing Platform
```

**Verification:**
- [ ] API URL configured correctly
- [ ] Environment variables set

#### 4.2 Build Frontend

```bash
cd frontend
pnpm install
pnpm build
```

**Verification:**
- [ ] Dependencies installed
- [ ] Build successful
- [ ] `dist/` directory created
- [ ] No build errors

#### 4.3 Deploy Frontend

**Option A: Static Hosting (Netlify, Vercel, etc.)**

```bash
# Deploy dist/ directory to your hosting provider
# Follow provider-specific instructions
```

**Option B: Nginx**

```bash
# Copy build to web root
sudo cp -r dist/* /var/www/uduxpass/

# Configure Nginx (see nginx.conf in repository)
sudo systemctl restart nginx
```

**Verification:**
- [ ] Frontend deployed
- [ ] Accessible via browser
- [ ] API calls working
- [ ] No console errors

---

### Step 5: Scanner App Deployment

#### 5.1 Configure Environment

Create `.env` file in `uduxpass-scanner-app/` directory:

```env
VITE_API_URL=https://api.your-domain.com
VITE_APP_NAME=uduXPass Scanner
```

**Verification:**
- [ ] API URL configured correctly
- [ ] Environment variables set

#### 5.2 Build Scanner App

```bash
cd uduxpass-scanner-app
pnpm install
pnpm build
```

**Verification:**
- [ ] Dependencies installed
- [ ] Build successful
- [ ] `client/dist/` directory created
- [ ] PWA manifest generated
- [ ] Service worker generated
- [ ] No build errors

#### 5.3 Deploy Scanner App

Deploy `client/dist/` directory to hosting provider.

**Verification:**
- [ ] Scanner app deployed
- [ ] Accessible via browser
- [ ] PWA installable
- [ ] API calls working
- [ ] QR scanner working
- [ ] Offline mode working

---

## Post-Deployment Verification

### 1. Backend API Tests

```bash
# Health check
curl https://api.your-domain.com/health

# User registration
curl -X POST https://api.your-domain.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Test",
    "lastName": "User",
    "email": "test@example.com",
    "phone": "+2348012345678",
    "password": "Test123!"
  }'

# User login
curl -X POST https://api.your-domain.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!"
  }'

# Get events
curl https://api.your-domain.com/api/events
```

**Verification:**
- [ ] Health check returns 200 OK
- [ ] Registration creates user
- [ ] Login returns JWT token
- [ ] Events endpoint returns data
- [ ] Response times < 200ms

### 2. Frontend Tests

**Manual Testing:**
- [ ] Homepage loads
- [ ] Events page displays events
- [ ] Event details page shows ticket tiers
- [ ] Registration form works
- [ ] Login form works
- [ ] Ticket selection works
- [ ] Checkout process works
- [ ] User profile accessible
- [ ] No console errors

### 3. Scanner App Tests

**Manual Testing:**
- [ ] Scanner login works
- [ ] Dashboard displays
- [ ] QR scanner opens
- [ ] Ticket validation works
- [ ] Session management works
- [ ] Offline mode works
- [ ] PWA installable

### 4. Integration Tests

**End-to-End Flow:**
- [ ] User can register
- [ ] User can login
- [ ] User can browse events
- [ ] User can select tickets
- [ ] User can checkout
- [ ] User receives confirmation email
- [ ] Scanner can validate tickets
- [ ] Scanner can mark tickets as used

---

## Security Checklist

### SSL/TLS
- [ ] SSL certificate installed
- [ ] HTTPS enabled for all domains
- [ ] HTTP redirects to HTTPS
- [ ] SSL certificate auto-renewal configured

### API Security
- [ ] CORS configured correctly
- [ ] Rate limiting enabled
- [ ] Input validation working
- [ ] SQL injection prevention verified
- [ ] XSS protection enabled
- [ ] CSRF protection enabled

### Authentication
- [ ] JWT tokens secure
- [ ] Password hashing working (bcrypt)
- [ ] Token expiration configured
- [ ] Refresh tokens implemented (if applicable)

### Database
- [ ] Database user has minimum required permissions
- [ ] SSL/TLS enabled for database connections
- [ ] Database backups configured
- [ ] Database access restricted to application servers

---

## Monitoring Setup

### Application Monitoring
- [ ] Error tracking configured (Sentry, etc.)
- [ ] Performance monitoring enabled
- [ ] Log aggregation configured
- [ ] Uptime monitoring enabled

### Alerts
- [ ] Error rate alerts configured
- [ ] Performance degradation alerts configured
- [ ] Uptime alerts configured
- [ ] Database alerts configured

---

## Backup and Recovery

### Database Backups
- [ ] Automated daily backups configured
- [ ] Backup retention policy set
- [ ] Backup restoration tested
- [ ] Off-site backup storage configured

### Application Backups
- [ ] Code repository backed up (GitHub)
- [ ] Environment configurations backed up
- [ ] SSL certificates backed up

---

## Performance Optimization

### Backend
- [ ] Database indexes optimized
- [ ] Query performance verified
- [ ] Connection pooling configured
- [ ] Caching enabled (Redis, if applicable)

### Frontend
- [ ] Static assets cached
- [ ] CDN configured
- [ ] Gzip compression enabled
- [ ] Images optimized
- [ ] Code splitting enabled

### Database
- [ ] Indexes created for frequently queried columns
- [ ] Query performance analyzed
- [ ] Connection pooling configured
- [ ] Vacuum and analyze scheduled

---

## Documentation

### Internal Documentation
- [ ] Deployment process documented
- [ ] Environment variables documented
- [ ] API endpoints documented
- [ ] Database schema documented
- [ ] Troubleshooting guide created

### User Documentation
- [ ] User guide created
- [ ] FAQ created
- [ ] Support contact information provided
- [ ] Terms of service published
- [ ] Privacy policy published

---

## Final Sign-Off

### Deployment Team
- [ ] Backend deployed and verified
- [ ] Frontend deployed and verified
- [ ] Scanner app deployed and verified
- [ ] Database configured and verified
- [ ] Security measures implemented
- [ ] Monitoring configured
- [ ] Backups configured
- [ ] Documentation complete

### Stakeholder Approval
- [ ] Technical lead approval
- [ ] Product owner approval
- [ ] Security team approval
- [ ] Operations team approval

---

## Rollback Plan

### If Deployment Fails:

1. **Backend Rollback:**
   ```bash
   # Stop new version
   sudo systemctl stop uduxpass-api
   
   # Restore previous version
   sudo systemctl start uduxpass-api-backup
   ```

2. **Frontend Rollback:**
   ```bash
   # Restore previous build
   sudo cp -r /var/www/uduxpass-backup/* /var/www/uduxpass/
   ```

3. **Database Rollback:**
   ```bash
   # Restore from backup
   pg_restore -U uduxpass -d uduxpass_db backup.dump
   ```

**Verification:**
- [ ] Previous version restored
- [ ] All services running
- [ ] Users can access platform
- [ ] No data loss

---

## Post-Deployment Tasks

### Week 1
- [ ] Monitor error rates
- [ ] Monitor performance metrics
- [ ] Gather user feedback
- [ ] Address critical issues

### Week 2-4
- [ ] Optimize performance based on metrics
- [ ] Address user feedback
- [ ] Plan feature enhancements
- [ ] Review security logs

---

## Support Contacts

**Technical Support:**
- Email: support@uduxpass.com
- Phone: +234 (0) 800 UDUXPASS

**Emergency Contacts:**
- DevOps Lead: [Contact Info]
- Backend Lead: [Contact Info]
- Frontend Lead: [Contact Info]

---

**Deployment Date:** _______________  
**Deployed By:** _______________  
**Verified By:** _______________  
**Approved By:** _______________

---

**Repository:** https://github.com/ArowuTest/uduXPass  
**Documentation:** See DOCKER_DEPLOYMENT_GUIDE.md and DEPLOYMENT_GUIDE.md  
**E2E Test Report:** See E2E_TEST_REPORT_FINAL.md

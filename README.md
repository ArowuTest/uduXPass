# uduXPass - Complete Event Ticketing Platform

**Version:** 1.0.0  
**Status:** Production Ready (with fixes needed)  
**Last Updated:** February 4, 2026

---

## 📦 Repository Structure

```
uduxpass-platform/
├── backend/              # Go API Server
│   ├── .env             # Backend environment variables
│   ├── cmd/             # Application entry points
│   ├── internal/        # Internal packages
│   ├── migrations/      # Database migrations
│   ├── uduxpass-api     # Compiled binary
│   └── start-as-postgres.sh
│
├── frontend/            # User-Facing React App
│   ├── src/            # React source code
│   ├── public/         # Static assets
│   ├── package.json    # Dependencies
│   └── Dockerfile      # Container config
│
├── scanner/            # Scanner PWA (React)
│   ├── src/           # Scanner app source
│   ├── public/        # PWA assets
│   ├── package.json   # Dependencies
│   └── vite.config.js # Build config
│
├── database/          # Database Files
│   ├── migrations/    # SQL migration files
│   │   ├── 001_initial_schema.sql
│   │   ├── 002_admin_users_schema.sql
│   │   ├── 003_scanner_system_schema.sql
│   │   └── 004_seed_data.sql
│   └── *.sql         # Utility SQL scripts
│
├── docs/             # Documentation
│   ├── test-qr-codes/     # Test QR code images
│   ├── CRITICAL_SCANNER_APP_BUGS.md
│   ├── UDUXPASS_FIXED_README.md
│   └── *.md          # All project documentation
│
├── docker-compose.yml     # Docker orchestration
└── README.md             # This file
```

---

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21+
- **Node.js** 18+
- **PostgreSQL** 15+
- **Redis** 7+
- **Docker** & Docker Compose (optional)

### 1. Start with Docker (Recommended)

```bash
# Clone/navigate to repository
cd uduxpass-platform

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

**Services will be available at:**
- Backend API: http://localhost:8080
- Frontend: http://localhost:5173
- Scanner App: http://localhost:3000
- PostgreSQL: localhost:5432
- Redis: localhost:6379

### 2. Manual Setup (Development)

#### A. Database Setup

```bash
# Start PostgreSQL
sudo systemctl start postgresql

# Create database
sudo -u postgres psql -c "CREATE DATABASE uduxpass;"

# Run migrations
cd database/migrations
sudo -u postgres psql uduxpass < 001_initial_schema.sql
sudo -u postgres psql uduxpass < 002_admin_users_schema.sql
sudo -u postgres psql uduxpass < 003_scanner_system_schema.sql
sudo -u postgres psql uduxpass < 004_seed_data.sql
```

#### B. Backend Setup

```bash
cd backend

# Configure environment
cp .env.example .env
# Edit .env with your settings

# Install dependencies
go mod download

# Build
go build -o uduxpass-api cmd/api/main.go

# Run
./start-as-postgres.sh
# OR
./uduxpass-api
```

**Backend runs on:** http://localhost:8080

#### C. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

**Frontend runs on:** http://localhost:5173

#### D. Scanner App Setup

```bash
cd scanner

# Install dependencies
npm install

# Start development server
PORT=3000 npm run dev
```

**Scanner runs on:** http://localhost:3000

---

## 🔐 Test Credentials

### Admin User
- **Email:** admin@uduxpass.com
- **Password:** Admin@123456

### Scanner Operators
- **Username:** scanner_lagos_1
- **Password:** Scanner@123
- **Name:** John Okafor
- **Location:** Lagos

### Regular Users
- **Email:** adeola.williams@gmail.com
- **Password:** User@123

---

## 🎫 Test Data

### Test Event
- **Name:** E2E Test Concert - Davido Live
- **ID:** ad48e795-dfa3-44a4-a8e0-7ddbdd8a689b
- **Date:** June 15, 2026
- **Venue:** Eko Atlantic Energy City, Lagos

### Test Tickets (with QR Codes)

1. **VIP Ticket**
   - Code: TICKET-VIP-001
   - QR: `UDUXPASS:ad48e795-dfa3-44a4-a8e0-7ddbdd8a689b:aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa`
   - Image: `docs/test-qr-codes/ticket-vip-001.png`

2. **Regular Ticket #1**
   - Code: TICKET-REG-001
   - QR: `UDUXPASS:ad48e795-dfa3-44a4-a8e0-7ddbdd8a689b:bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb`
   - Image: `docs/test-qr-codes/ticket-reg-001.png`

3. **Regular Ticket #2**
   - Code: TICKET-REG-002
   - QR: `UDUXPASS:ad48e795-dfa3-44a4-a8e0-7ddbdd8a689b:cccccccc-cccc-cccc-cccc-cccccccccccc`
   - Image: `docs/test-qr-codes/ticket-reg-002.png`

---

## 🧪 Testing

### Backend API Tests

```bash
cd backend

# Health check
curl http://localhost:8080/health

# Admin login
curl -X POST http://localhost:8080/v1/auth/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@uduxpass.com","password":"Admin@123456"}'

# Scanner login
curl -X POST http://localhost:8080/v1/scanner/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"scanner_lagos_1","password":"Scanner@123"}'
```

### Scanner App E2E Test

1. Open http://localhost:3000
2. Login with scanner credentials
3. Select event from dashboard
4. Use test QR codes from `docs/test-qr-codes/`
5. Verify ticket validation

---

## 🐛 Known Issues

### Critical Issues (Must Fix)

1. **Scanner App: Event Data Not Displaying**
   - **Problem:** Dashboard shows "Invalid Date" instead of event details
   - **Root Cause:** Data mapping mismatch between API and frontend
   - **Status:** Documented in `docs/CRITICAL_SCANNER_APP_BUGS.md`
   - **Priority:** P0

2. **Scanner App: Date Formatting Broken**
   - **Problem:** All dates show as "Invalid Date"
   - **Root Cause:** Date parsing not handling ISO 8601 format
   - **Status:** Documented
   - **Priority:** P0

### Medium Priority Issues

3. **Scanner App: Event Card Not Clickable**
   - **Status:** Needs implementation
   - **Priority:** P1

4. **Scanner App: Missing Event Images**
   - **Status:** Needs implementation
   - **Priority:** P1

**See:** `docs/CRITICAL_SCANNER_APP_BUGS.md` for complete details and fixes.

---

## 📝 Environment Variables

### Backend (.env)

```env
# Server
PORT=8080
HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=uduxpass
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Payment
PAYSTACK_SECRET_KEY=your_paystack_secret
PAYSTACK_PUBLIC_KEY=your_paystack_public

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
```

### Frontend (.env)

```env
VITE_API_URL=http://localhost:8080/v1
VITE_APP_NAME=uduXPass
```

### Scanner (.env)

```env
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=uduXPass Scanner
```

---

## 🏗️ Architecture

### Technology Stack

**Backend:**
- Go 1.21
- Gin (HTTP framework)
- sqlx (Database)
- JWT authentication
- bcrypt password hashing

**Frontend:**
- React 18
- TypeScript
- Vite
- TailwindCSS
- Axios

**Scanner:**
- React 18
- TypeScript
- Vite
- html5-qrcode
- PWA capabilities

**Database:**
- PostgreSQL 15
- Redis 7

**Infrastructure:**
- Docker
- Docker Compose
- Nginx (production)

---

## 📚 API Documentation

### Base URL
```
http://localhost:8080/v1
```

### Authentication Endpoints

#### Admin Login
```http
POST /auth/admin/login
Content-Type: application/json

{
  "email": "admin@uduxpass.com",
  "password": "Admin@123456"
}
```

#### Scanner Login
```http
POST /scanner/auth/login
Content-Type: application/json

{
  "username": "scanner_lagos_1",
  "password": "Scanner@123"
}
```

#### User Registration
```http
POST /auth/email/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "Password@123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+2348012345678"
}
```

### Scanner Endpoints

#### Get Assigned Events
```http
GET /scanner/events
Authorization: Bearer {token}
```

#### Validate Ticket
```http
POST /scanner/validate
Authorization: Bearer {token}
Content-Type: application/json

{
  "ticket_code": "TICKET-VIP-001",
  "event_id": "ad48e795-dfa3-44a4-a8e0-7ddbdd8a689b"
}
```

---

## 🔧 Development

### Backend Development

```bash
cd backend

# Run with hot reload (if air is installed)
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint
golangci-lint run
```

### Frontend Development

```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint
npm run lint
```

### Scanner Development

```bash
cd scanner

# Development server
npm run dev

# Build for production
npm run build

# Test PWA
npm run preview
```

---

## 📦 Deployment

### Production Build

```bash
# Build backend
cd backend
go build -o uduxpass-api cmd/api/main.go

# Build frontend
cd frontend
npm run build

# Build scanner
cd scanner
npm run build
```

### Docker Deployment

```bash
# Build all images
docker-compose build

# Start in production mode
docker-compose -f docker-compose.prod.yml up -d
```

---

## 🤝 Contributing

1. Create feature branch
2. Make changes
3. Test thoroughly
4. Submit pull request

---

## 📄 License

Proprietary - All rights reserved

---

## 📞 Support

For issues and questions:
- Check `docs/` directory for detailed documentation
- Review `docs/CRITICAL_SCANNER_APP_BUGS.md` for known issues
- Contact development team

---

## ✅ Production Readiness Checklist

- [x] Backend API running
- [x] Database schema created
- [x] Seed data loaded
- [x] Authentication working
- [x] Frontend deployed
- [x] Scanner app deployed
- [ ] Critical bugs fixed (see CRITICAL_SCANNER_APP_BUGS.md)
- [ ] E2E testing completed
- [ ] Performance testing
- [ ] Security audit
- [ ] Load testing
- [ ] Documentation complete

**Current Status:** 85% Production Ready

**Blocking Issues:** 2 critical scanner app bugs

---

**Built with ❤️ for event organizers**

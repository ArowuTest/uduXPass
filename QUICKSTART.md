# uduXPass Platform - Quick Start Guide

## 🚀 Get Started in 5 Minutes

### Option 1: One-Command Start (Easiest)

```bash
cd uduxpass-platform
./start-all.sh
```

That's it! All services will start automatically.

**Access:**
- Backend: http://localhost:8080
- Frontend: http://localhost:5173  
- Scanner: http://localhost:3000

**Stop everything:**
```bash
./stop-all.sh
```

---

### Option 2: Docker (Recommended for Production)

```bash
cd uduxpass-platform
docker-compose up -d
```

**Stop:**
```bash
docker-compose down
```

---

### Option 3: Manual Start (For Development)

#### 1. Start Backend
```bash
cd backend
./start-as-postgres.sh
```

#### 2. Start Frontend (new terminal)
```bash
cd frontend
npm install  # first time only
npm run dev
```

#### 3. Start Scanner (new terminal)
```bash
cd scanner
npm install  # first time only
PORT=3000 npm run dev
```

---

## 🔐 Login Credentials

### Admin Portal
- Email: `admin@uduxpass.com`
- Password: `Admin@123456`

### Scanner App
- Username: `scanner_lagos_1`
- Password: `Scanner@123`

### User Account
- Email: `adeola.williams@gmail.com`
- Password: `User@123`

---

## 🧪 Test the Scanner

1. Open http://localhost:3000
2. Login with scanner credentials
3. Select "E2E Test Concert - Davido Live"
4. Use test QR codes from `docs/test-qr-codes/`
5. Scan and validate tickets!

---

## 📦 What's Included

- ✅ Complete backend API (Go)
- ✅ User-facing frontend (React)
- ✅ Scanner PWA app (React)
- ✅ PostgreSQL database with seed data
- ✅ Test QR codes
- ✅ Complete documentation
- ✅ Docker setup
- ✅ Startup scripts

---

## 🐛 Known Issues

**2 Critical Bugs in Scanner App:**
1. Event data displays as "Invalid Date"
2. Date formatting broken

See `docs/CRITICAL_SCANNER_APP_BUGS.md` for details and fixes.

---

## 📚 Full Documentation

- **README.md** - Complete platform documentation
- **docs/UDUXPASS_FIXED_README.md** - Detailed setup guide
- **docs/CRITICAL_SCANNER_APP_BUGS.md** - Bug reports and fixes
- **docker-compose.yml** - Container orchestration

---

## 💡 Tips

**View logs:**
```bash
tail -f logs/backend.log
tail -f logs/frontend.log
tail -f logs/scanner.log
```

**Check health:**
```bash
curl http://localhost:8080/health
```

**Reset database:**
```bash
cd database/migrations
sudo -u postgres psql uduxpass < 001_initial_schema.sql
sudo -u postgres psql uduxpass < 004_seed_data.sql
```

---

## 🆘 Troubleshooting

**Port already in use:**
```bash
./stop-all.sh
./start-all.sh
```

**Backend not connecting to database:**
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check backend .env file has correct DB settings
cat backend/.env
```

**Frontend/Scanner can't reach backend:**
```bash
# Check backend is running
curl http://localhost:8080/health

# Check CORS settings in backend/.env
```

---

## ✅ Next Steps

1. ✅ Start all services
2. ✅ Login to each app
3. ✅ Test scanner with QR codes
4. 🔧 Fix critical bugs (see docs/)
5. 🧪 Run comprehensive tests
6. 🚀 Deploy to production

---

**Need help?** Check the full README.md or docs/ directory.

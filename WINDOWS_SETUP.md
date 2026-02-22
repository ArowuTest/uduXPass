# uduXPass Platform - Windows Setup Guide

## 📋 Prerequisites

Before starting, install these tools:

### 1. PostgreSQL 15+
- Download: https://www.postgresql.org/download/windows/
- During installation, remember your postgres password
- Default port: 5432

### 2. Go 1.21+
- Download: https://golang.org/dl/
- Install and verify: `go version`

### 3. Node.js 18+
- Download: https://nodejs.org/
- Install and verify: `node --version`

### 4. Git (Optional)
- Download: https://git-scm.com/download/win

---

## 🚀 Quick Start

### Step 1: Extract the ZIP
Extract `uduxpass-platform.zip` to a folder like:
```
C:\uduxpass-platform\
```

### Step 2: Setup Database

Open **Command Prompt** or **PowerShell** as Administrator:

```cmd
cd C:\uduxpass-platform\database\migrations

REM Create database
psql -U postgres -c "CREATE DATABASE uduxpass;"

REM Run migrations
psql -U postgres -d uduxpass -f 001_initial_schema.sql
psql -U postgres -d uduxpass -f 002_admin_users_schema.sql
psql -U postgres -d uduxpass -f 003_scanner_system_schema.sql
psql -U postgres -d uduxpass -f 004_seed_data.sql
```

### Step 3: Configure Backend

Edit `backend\.env` file with your settings:

```env
# Database (update with your PostgreSQL password)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=YOUR_POSTGRES_PASSWORD
DB_NAME=uduxpass
DB_SSLMODE=disable

# Server
PORT=8080
HOST=0.0.0.0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Step 4: Install Dependencies

#### Backend
```cmd
cd backend
go mod download
go build -o uduxpass-api.exe cmd\api\main.go
```

#### Frontend
```cmd
cd frontend
npm install
```

#### Scanner
```cmd
cd scanner
npm install
```

### Step 5: Start All Services

**Option A: Use Batch Script (Easiest)**
```cmd
cd C:\uduxpass-platform
start-all.bat
```

**Option B: Manual Start**

Open 3 separate Command Prompt windows:

**Window 1 - Backend:**
```cmd
cd C:\uduxpass-platform\backend
uduxpass-api.exe
```

**Window 2 - Frontend:**
```cmd
cd C:\uduxpass-platform\frontend
npm run dev
```

**Window 3 - Scanner:**
```cmd
cd C:\uduxpass-platform\scanner
set PORT=3000
npm run dev
```

---

## 🌐 Access the Platform

After starting all services:

- **Backend API:** http://localhost:8080
- **Frontend:** http://localhost:5173
- **Scanner App:** http://localhost:3000

---

## 🔐 Test Credentials

### Admin Portal (Frontend)
- Email: `admin@uduxpass.com`
- Password: `Admin@123456`

### Scanner App
- Username: `scanner_lagos_1`
- Password: `Scanner@123`

### Regular User
- Email: `adeola.williams@gmail.com`
- Password: `User@123`

---

## 🧪 Test the Scanner

1. Open http://localhost:3000
2. Login with scanner credentials
3. You should see "E2E Test Concert - Davido Live"
4. Use test QR codes from `docs\test-qr-codes\`
5. Scan and validate tickets!

---

## 🐛 Troubleshooting

### PostgreSQL Connection Error

**Error:** `connection refused` or `password authentication failed`

**Solution:**
1. Check PostgreSQL is running:
   - Open Services (Win + R, type `services.msc`)
   - Find "postgresql-x64-15" service
   - Ensure it's running

2. Verify password in `backend\.env` matches your PostgreSQL password

3. Test connection:
   ```cmd
   psql -U postgres -d uduxpass
   ```

### Port Already in Use

**Error:** `port 8080 already in use`

**Solution:**
```cmd
REM Find process using port
netstat -ano | findstr :8080

REM Kill process (replace PID with actual process ID)
taskkill /PID <PID> /F
```

### Go Build Error

**Error:** `go: command not found`

**Solution:**
1. Ensure Go is installed
2. Add Go to PATH:
   - Right-click "This PC" → Properties
   - Advanced system settings → Environment Variables
   - Add `C:\Program Files\Go\bin` to PATH

### Node/NPM Error

**Error:** `npm: command not found`

**Solution:**
1. Ensure Node.js is installed
2. Restart Command Prompt after installation
3. Verify: `node --version` and `npm --version`

### Frontend/Scanner Can't Reach Backend

**Solution:**
1. Check backend is running: http://localhost:8080/health
2. Check CORS settings in `backend\.env`
3. Ensure firewall allows connections to port 8080

---

## 📁 Project Structure

```
C:\uduxpass-platform\
├── backend\              # Go API Server
│   ├── .env             # Backend configuration
│   ├── cmd\             # Application entry
│   ├── internal\        # Internal packages
│   └── uduxpass-api.exe # Compiled binary
│
├── frontend\            # User-Facing React App
│   ├── src\            # React source
│   ├── .env            # Frontend config
│   └── package.json    # Dependencies
│
├── scanner\            # Scanner PWA
│   ├── src\           # Scanner source
│   ├── .env           # Scanner config
│   └── package.json   # Dependencies
│
├── database\          # Database Files
│   └── migrations\    # SQL migrations
│
├── docs\              # Documentation
│   └── test-qr-codes\ # Test QR images
│
├── start-all.bat      # Windows startup script
├── README.md          # Main documentation
└── QUICKSTART.md      # Quick start guide
```

---

## 🔄 Updating the Platform

### Update Backend
```cmd
cd backend
go build -o uduxpass-api.exe cmd\api\main.go
```

### Update Frontend/Scanner
```cmd
cd frontend
npm install
npm run build
```

---

## 🛑 Stopping Services

Press `Ctrl + C` in each Command Prompt window running a service.

Or close the windows opened by `start-all.bat`.

---

## 📚 Additional Resources

- **Full Documentation:** `README.md`
- **Quick Start:** `QUICKSTART.md`
- **Bug Reports:** `docs\CRITICAL_SCANNER_APP_BUGS.md`
- **API Documentation:** See README.md

---

## 💡 Development Tips

### Hot Reload

Both frontend and scanner support hot reload - changes to code will automatically refresh the browser.

### Database Reset

To reset database to initial state:
```cmd
cd database\migrations
psql -U postgres -d uduxpass -f 001_initial_schema.sql
psql -U postgres -d uduxpass -f 004_seed_data.sql
```

### View Logs

Backend logs appear in the Command Prompt window.

Frontend/Scanner logs appear in browser console (F12).

---

## ✅ Production Deployment

For production deployment on Windows Server:

1. Build all components
2. Use Windows Service to run backend
3. Use IIS or Nginx for frontend/scanner
4. Configure SSL certificates
5. Set up proper firewall rules
6. Use production database credentials

See `README.md` for detailed deployment instructions.

---

## 🆘 Getting Help

If you encounter issues:

1. Check this guide's Troubleshooting section
2. Review `docs\CRITICAL_SCANNER_APP_BUGS.md`
3. Check backend logs in Command Prompt
4. Check browser console (F12) for frontend errors

---

**Built for Windows with ❤️**

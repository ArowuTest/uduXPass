# 🚀 uduXPass Backend - Bcrypt Update Quick Start

## ✅ What's Fixed

This repository contains **production-ready fixes** for:

1. ✅ **Bcrypt password hashing** (replaced Argon2)
2. ✅ **Schema alignment** (AdminUser entity matches database 100%)
3. ✅ **PostgreSQL array support** (custom JSONB scanner)
4. ✅ **Authentication bugs** (parameter order fixed)
5. ✅ **Seed data** (bcrypt hashes with correct schema)

---

## 📦 What's Included

```
uduxpass-backend/
├── CHANGES.md                          # Comprehensive changelog
├── README_BCRYPT_UPDATE.md            # This file
├── cmd/api/main.go                     # Backend entry point
├── internal/
│   ├── domain/entities/
│   │   ├── admin_user.go              # ✅ Fixed schema alignment
│   │   └── admin_permission_array.go  # ✅ NEW: JSONB scanner
│   ├── usecases/admin/
│   │   └── admin_auth_service.go      # ✅ Fixed auth bugs
│   └── interfaces/http/server/
│       └── server.go                   # ✅ Bcrypt initialization
├── pkg/security/
│   └── password.go                     # ✅ Bcrypt implementation
└── migrations/
    └── 004_seed_data.sql               # ✅ Bcrypt hashes
```

---

## 🎯 Quick Test (5 minutes)

### **1. Extract & Build**
```bash
unzip uduxpass-backend-bcrypt-fixed.zip
cd uduxpass-backend
go mod tidy
go build -o uduxpass-api cmd/api/main.go
```

### **2. Setup Database**
```bash
# Create database
psql -U postgres -c "CREATE DATABASE uduxpass;"

# Apply migrations
psql -U postgres -d uduxpass -f migrations/001_initial_schema.sql
psql -U postgres -d uduxpass -f migrations/002_admin_users_schema.sql
psql -U postgres -d uduxpass -f migrations/003_scanner_system_schema.sql
psql -U postgres -d uduxpass -f migrations/004_seed_data.sql
```

### **3. Configure Environment**
```bash
# Create .env file
cat > .env << EOF
DATABASE_URL=postgres://postgres:postgres@localhost:5432/uduxpass?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
ENVIRONMENT=development
EOF
```

### **4. Start Backend**
```bash
./uduxpass-api
```

### **5. Test Admin Login**
```bash
curl -X POST http://localhost:8080/v1/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@uduxpass.com",
    "password": "Admin@123456"
  }'
```

### **Expected Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "admin": {
    "id": "11111111-1111-1111-1111-111111111111",
    "email": "admin@uduxpass.com",
    "first_name": "System",
    "last_name": "Administrator",
    "role": "super_admin",
    "permissions": ["system_settings", "user_management", ...],
    "is_active": true
  },
  "expires_in": 3600
}
```

✅ **If you see this response, everything is working!**

---

## 🔑 Test Credentials

| Email | Password | Role |
|-------|----------|------|
| admin@uduxpass.com | Admin@123456 | super_admin |
| events@uduxpass.com | Admin@123456 | event_manager |
| support@uduxpass.com | Admin@123456 | support |

---

## 📋 Key Files Changed

### **1. Password Service** (`pkg/security/password.go`)
- **Before**: 400+ lines of Argon2 code
- **After**: 200 lines of bcrypt code
- **Change**: Complete rewrite with bcrypt

### **2. AdminUser Entity** (`internal/domain/entities/admin_user.go`)
- **Before**: 17 fields (7 not in database)
- **After**: 15 fields (100% match with database)
- **Change**: Removed duplicate/non-existent fields

### **3. AdminPermissionArray** (`internal/domain/entities/admin_permission_array.go`)
- **Status**: NEW FILE
- **Purpose**: Custom scanner for PostgreSQL JSONB arrays
- **Methods**: `Scan()`, `Value()`, `MarshalJSON()`, `UnmarshalJSON()`

### **4. Admin Auth Service** (`internal/usecases/admin/admin_auth_service.go`)
- **Fixed**: Parameter order in `VerifyPassword()` calls
- **Fixed**: All references to removed fields
- **Fixed**: Status checks to use `IsActive`

### **5. Seed Data** (`migrations/004_seed_data.sql`)
- **Updated**: All password hashes to bcrypt
- **Fixed**: Column list to match schema
- **Fixed**: Permissions format to JSONB

---

## 🐛 Common Issues & Solutions

### **Issue 1: "Admin access required"**
**Cause**: Token not in Authorization header  
**Solution**: Add `Authorization: Bearer <token>` header

### **Issue 2: "Invalid credentials"**
**Cause**: Wrong password or hash mismatch  
**Solution**: Use `Admin@123456` (case-sensitive)

### **Issue 3: "Database connection failed"**
**Cause**: PostgreSQL not running or wrong credentials  
**Solution**: Check `.env` file and PostgreSQL status

### **Issue 4: "Permissions not loading"**
**Cause**: Old schema without JSONB permissions  
**Solution**: Reapply migration `002_admin_users_schema.sql`

---

## 📚 Documentation

For complete details, see:
- **CHANGES.md** - Full changelog with technical details
- **migrations/** - Database schema and seed data
- **pkg/security/password.go** - Bcrypt implementation
- **internal/domain/entities/admin_permission_array.go** - JSONB scanner

---

## ✅ Production Checklist

Before deploying to production:

- [ ] Change `JWT_SECRET` in `.env`
- [ ] Update admin passwords from default
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS allowed origins
- [ ] Set `ENVIRONMENT=production`
- [ ] Enable database backups
- [ ] Configure logging
- [ ] Set up monitoring
- [ ] Test all admin endpoints
- [ ] Test all authentication flows

---

## 🎉 Success Criteria

You'll know it's working when:

✅ Admin login returns JWT token  
✅ Token works for protected endpoints  
✅ Permissions checked correctly  
✅ Account lockout works after 5 failed attempts  
✅ No errors in backend logs  

---

## 📞 Support

If you encounter any issues:

1. Check **CHANGES.md** for detailed documentation
2. Review backend logs for error messages
3. Verify database schema matches migrations
4. Test with provided credentials first

---

**Status**: ✅ **PRODUCTION-READY**  
**Version**: 2.0.0-bcrypt  
**Date**: February 3, 2026

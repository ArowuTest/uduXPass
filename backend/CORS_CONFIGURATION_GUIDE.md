# CORS Configuration Guide for uduXPass Backend

## Overview

The uduXPass backend now supports flexible CORS configuration for all deployment scenarios:
- ✅ Local development
- ✅ Manus sandbox testing
- ✅ Vercel/Netlify deployment
- ✅ Render/Railway deployment
- ✅ Custom domains
- ✅ Production environments

## Configuration Methods

### Method 1: Environment Variable (Recommended for Production)

Set the `CORS_ALLOWED_ORIGINS` environment variable with a comma-separated list of allowed origins:

```bash
# Production example with specific domains
CORS_ALLOWED_ORIGINS="https://uduxpass.com,https://www.uduxpass.com,https://admin.uduxpass.com"

# Vercel deployment example
CORS_ALLOWED_ORIGINS="https://uduxpass.vercel.app,https://uduxpass-admin.vercel.app"

# Render deployment example
CORS_ALLOWED_ORIGINS="https://uduxpass.onrender.com,https://uduxpass-admin.onrender.com"

# Netlify deployment example
CORS_ALLOWED_ORIGINS="https://uduxpass.netlify.app,https://admin-uduxpass.netlify.app"

# Multiple deployment environments
CORS_ALLOWED_ORIGINS="https://uduxpass.com,https://staging.uduxpass.com,https://uduxpass.vercel.app"
```

### Method 2: Allow All Origins (Development Only)

For development/testing, you can allow all origins:

```bash
CORS_ALLOWED_ORIGINS="*"
```

**⚠️ WARNING:** Never use `*` in production! This is a security risk.

### Method 3: Environment-Based Defaults

If you don't set `CORS_ALLOWED_ORIGINS`, the backend uses smart defaults based on the `ENVIRONMENT` variable:

```bash
# Development/Sandbox - Allows all origins
ENVIRONMENT="development"  # or "sandbox"

# Production - Only allows localhost (you should set CORS_ALLOWED_ORIGINS)
ENVIRONMENT="production"
```

## Deployment Examples

### Local Development

No configuration needed! The backend automatically allows:
- `http://localhost:3000` (Scanner app)
- `http://localhost:5173` (Frontend)
- `http://localhost:8080` (Backend)

Just run:
```bash
cd backend
DATABASE_URL="postgres://user:pass@localhost:5432/uduxpass" go run cmd/api/main.go
```

### Vercel Deployment

1. Deploy frontend to Vercel:
   ```bash
   cd frontend
   vercel deploy --prod
   # Output: https://uduxpass.vercel.app
   ```

2. Set backend environment variable:
   ```bash
   CORS_ALLOWED_ORIGINS="https://uduxpass.vercel.app,https://uduxpass-*.vercel.app"
   ```

3. Deploy backend to Render/Railway with the environment variable

### Render Deployment

**Frontend on Render:**
1. Create new Static Site
2. Build command: `npm run build`
3. Publish directory: `dist`
4. Get URL: `https://uduxpass.onrender.com`

**Backend on Render:**
1. Create new Web Service
2. Add environment variable:
   ```
   CORS_ALLOWED_ORIGINS=https://uduxpass.onrender.com
   ```
3. Deploy

### Railway Deployment

**Frontend on Railway:**
1. Deploy from GitHub
2. Get URL: `https://uduxpass.up.railway.app`

**Backend on Railway:**
1. Deploy from GitHub
2. Add environment variable:
   ```
   CORS_ALLOWED_ORIGINS=https://uduxpass.up.railway.app
   ```

### Netlify Deployment

**Frontend on Netlify:**
1. Deploy from GitHub
2. Build command: `npm run build`
3. Publish directory: `dist`
4. Get URL: `https://uduxpass.netlify.app`

**Backend (deploy elsewhere):**
```bash
CORS_ALLOWED_ORIGINS="https://uduxpass.netlify.app"
```

### Custom Domain

If you have a custom domain:

```bash
# Single domain
CORS_ALLOWED_ORIGINS="https://uduxpass.com"

# Multiple subdomains
CORS_ALLOWED_ORIGINS="https://uduxpass.com,https://www.uduxpass.com,https://admin.uduxpass.com,https://scanner.uduxpass.com"
```

## Complete Environment Variables

Here's a complete example for production deployment:

```bash
# Database
DATABASE_URL="postgres://user:pass@host:5432/uduxpass?sslmode=require"

# Server
PORT="8080"
ENVIRONMENT="production"

# CORS - Add your frontend domains here
CORS_ALLOWED_ORIGINS="https://uduxpass.com,https://www.uduxpass.com,https://admin.uduxpass.com"

# JWT
JWT_SECRET="your-super-secret-jwt-key-change-this-in-production"

# Payment Providers
PAYSTACK_SECRET_KEY="sk_live_your_paystack_key"
MOMO_API_KEY="your_momo_api_key"

# Email (SendGrid example)
SMTP_HOST="smtp.sendgrid.net"
SMTP_PORT="587"
SMTP_USERNAME="apikey"
SMTP_PASSWORD="your_sendgrid_api_key"
SMTP_FROM_EMAIL="noreply@uduxpass.com"
SMTP_FROM_NAME="uduXPass"
```

## Testing CORS Configuration

### Test 1: Check Allowed Origins

```bash
curl -I -X OPTIONS http://localhost:8080/v1/events \
  -H "Origin: https://uduxpass.com" \
  -H "Access-Control-Request-Method: GET"
```

Expected response should include:
```
Access-Control-Allow-Origin: https://uduxpass.com
Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Credentials: true
```

### Test 2: Test from Browser Console

Open your frontend in browser and run:

```javascript
fetch('https://api.uduxpass.com/v1/events')
  .then(res => res.json())
  .then(data => console.log('Success:', data))
  .catch(err => console.error('CORS Error:', err));
```

If you see a CORS error, check:
1. Is your frontend domain in `CORS_ALLOWED_ORIGINS`?
2. Did you restart the backend after changing the environment variable?
3. Is the protocol correct (http vs https)?

## Common Issues & Solutions

### Issue 1: "CORS policy: No 'Access-Control-Allow-Origin' header"

**Solution:** Add your frontend domain to `CORS_ALLOWED_ORIGINS`:
```bash
CORS_ALLOWED_ORIGINS="https://your-frontend-domain.com"
```

### Issue 2: "CORS policy: The value of the 'Access-Control-Allow-Credentials' header"

**Solution:** Make sure your frontend is sending credentials:
```javascript
fetch(url, {
  credentials: 'include'  // Add this
})
```

### Issue 3: Works on localhost but not on deployed domain

**Solution:** You need to add the deployed domain to CORS:
```bash
# Add both localhost AND deployed domain
CORS_ALLOWED_ORIGINS="http://localhost:5173,https://uduxpass.vercel.app"
```

### Issue 4: Wildcard subdomain not working

**Solution:** CORS doesn't support wildcard subdomains. List each subdomain explicitly:
```bash
# Instead of *.uduxpass.com, use:
CORS_ALLOWED_ORIGINS="https://app.uduxpass.com,https://admin.uduxpass.com,https://scanner.uduxpass.com"
```

## Security Best Practices

### ✅ DO:
- Use specific domains in production
- Use HTTPS in production
- List all subdomains explicitly
- Use environment variables for configuration
- Test CORS after deployment

### ❌ DON'T:
- Use `CORS_ALLOWED_ORIGINS="*"` in production
- Allow `http://` origins in production (use `https://`)
- Hardcode origins in the code
- Forget to restart backend after changing environment variables

## Migration Guide

### From AllowAllOrigins to Specific Origins

If you're currently using `AllowAllOrigins = true`, migrate to specific origins:

**Before:**
```go
corsConfig.AllowAllOrigins = true
```

**After:**
```bash
# Set environment variable
CORS_ALLOWED_ORIGINS="https://uduxpass.com,https://admin.uduxpass.com"
```

### From Hardcoded Origins to Environment Variable

**Before:**
```go
corsConfig.AllowOrigins = []string{
    "http://localhost:3000",
    "https://uduxpass.com",
}
```

**After:**
```bash
# Set environment variable
CORS_ALLOWED_ORIGINS="http://localhost:3000,https://uduxpass.com"
```

## Support

If you encounter CORS issues:

1. Check backend logs for CORS configuration on startup
2. Verify environment variable is set correctly
3. Test with curl to isolate the issue
4. Check browser console for detailed CORS error messages
5. Ensure backend was restarted after changing environment variables

## Summary

The uduXPass backend CORS configuration is now:
- ✅ **Flexible** - Works with any deployment platform
- ✅ **Secure** - Environment-based defaults prevent accidental exposure
- ✅ **Simple** - Just set one environment variable
- ✅ **Production-ready** - Supports multiple domains and subdomains

For most deployments, you just need to set:
```bash
CORS_ALLOWED_ORIGINS="https://your-frontend-domain.com"
```

That's it! 🎉

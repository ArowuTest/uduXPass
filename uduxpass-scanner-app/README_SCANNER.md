# uduXPass Scanner App

Professional Progressive Web App (PWA) for scanning and validating event tickets with QR code support.

## Features

✅ **Scanner Authentication** - Secure login for event scanners  
✅ **QR Code Scanning** - Real-time camera-based ticket scanning  
✅ **Session Management** - Create and manage scanning sessions  
✅ **Ticket Validation** - Instant validation with clear success/error states  
✅ **Session History** - View past scanning sessions and statistics  
✅ **Mobile-First Design** - Optimized for handheld scanning devices  
✅ **Offline Support** - PWA capabilities for reliable operation  

## Technology Stack

- **React 19** - Modern UI framework
- **TypeScript** - Type-safe development
- **TailwindCSS 4** - Utility-first styling
- **html5-qrcode** - QR code scanning
- **Wouter** - Lightweight routing
- **shadcn/ui** - High-quality UI components
- **Vite** - Fast build tool

## Design Philosophy

**Professional Event Tech** - The scanner app follows a "Professional Event Tech" design philosophy:

- **Clarity Over Decoration** - Every element serves a clear purpose
- **Speed & Efficiency** - Optimized for rapid scanning workflows
- **Trust & Reliability** - Professional aesthetics conveying security
- **Mobile-First Precision** - Designed for handheld devices in event environments

### Color Scheme

- **Primary Blue** (#1E40AF) - Trust, professionalism, branding
- **Success Green** (#10B981) - Valid tickets, positive feedback
- **Error Red** (#EF4444) - Invalid tickets, clear alerts
- **Clean Neutrals** - White backgrounds, subtle grays

## Setup Instructions

### 1. Install Dependencies

```bash
pnpm install
```

### 2. Configure Backend API

The app needs to connect to the uduXPass backend. Set the backend API URL:

**Option A: Environment Variable**
Create a `.env` file in the project root:

```env
VITE_API_BASE_URL=http://your-backend-url:8080/api/v1
```

**Option B: Direct Configuration**
Edit `client/src/lib/api.ts` and update the `API_BASE_URL`:

```typescript
const API_BASE_URL = 'http://your-backend-url:8080/api/v1';
```

### 3. Run Development Server

```bash
pnpm dev
```

The app will be available at `http://localhost:3000`

### 4. Build for Production

```bash
pnpm build
```

The production build will be in the `dist/` directory.

## Usage Guide

### Scanner Login

1. Open the app and navigate to the login page
2. Enter scanner credentials (email and password)
3. Click "Login to Scanner"

**Test Credentials** (from backend seed data):
- Email: `scanner1`
- Password: `Scanner@123!`

### Create Scanning Session

1. From the dashboard, click "Start New Session"
2. Select an event from the dropdown
3. Enter the location/entrance (e.g., "Main Entrance")
4. Optionally add notes
5. Click "Start Scanning Session"

### Scan Tickets

1. From the dashboard, click "Scan Ticket"
2. Point the camera at the ticket's QR code
3. The app will automatically scan and validate
4. View the result (green for valid, red for invalid)
5. Click "Scan Next Ticket" to continue

### View Session History

1. From the dashboard, click "Session History"
2. View all past and active scanning sessions
3. See statistics for each session

## Project Structure

```
client/
├── src/
│   ├── components/
│   │   └── ui/          # shadcn/ui components
│   ├── contexts/
│   │   ├── AuthContext.tsx    # Authentication state
│   │   └── ThemeContext.tsx   # Theme management
│   ├── lib/
│   │   ├── api.ts             # Backend API client
│   │   └── utils.ts           # Utility functions
│   ├── pages/
│   │   ├── Login.tsx          # Scanner login
│   │   ├── Dashboard.tsx      # Main dashboard
│   │   ├── Scanner.tsx        # QR scanner
│   │   ├── ValidationSuccess.tsx  # Valid ticket
│   │   ├── ValidationError.tsx    # Invalid ticket
│   │   ├── CreateSession.tsx  # Create session
│   │   └── SessionHistory.tsx # Session history
│   ├── App.tsx          # Routes and providers
│   ├── main.tsx         # Entry point
│   └── index.css        # Global styles
└── index.html           # HTML template
```

## API Integration

The app integrates with the uduXPass backend APIs:

- `POST /api/v1/scanner/auth/login` - Scanner authentication
- `GET /api/v1/scanner/events` - List available events
- `POST /api/v1/scanner/sessions` - Create scanning session
- `GET /api/v1/scanner/sessions` - List sessions
- `POST /api/v1/scanner/validate` - Validate ticket
- `GET /api/v1/scanner/sessions/:id/stats` - Session statistics

## Camera Permissions

The scanner requires camera access to scan QR codes. When prompted:

1. **iOS Safari**: Tap "Allow" when prompted for camera access
2. **Android Chrome**: Tap "Allow" when prompted for camera access
3. **Desktop**: Click "Allow" in the browser permission dialog

If camera access is denied, the app provides a manual entry fallback.

## PWA Installation

The app can be installed as a Progressive Web App:

### iOS
1. Open the app in Safari
2. Tap the Share button
3. Tap "Add to Home Screen"
4. Tap "Add"

### Android
1. Open the app in Chrome
2. Tap the menu (three dots)
3. Tap "Add to Home Screen"
4. Tap "Add"

### Desktop
1. Open the app in Chrome
2. Click the install icon in the address bar
3. Click "Install"

## Troubleshooting

### Camera Not Working

- **Check permissions**: Ensure camera permissions are granted
- **HTTPS required**: Camera access requires HTTPS (except localhost)
- **Try manual entry**: Use the "Enter Ticket Code Manually" fallback

### Login Failed

- **Check backend**: Ensure the backend API is running
- **Check URL**: Verify `VITE_API_BASE_URL` is correct
- **Check credentials**: Use valid scanner credentials

### No Events Available

- **Check backend**: Ensure events exist in the database
- **Check permissions**: Scanner must have access to events

## Development

### Adding New Features

1. Create new components in `client/src/components/`
2. Create new pages in `client/src/pages/`
3. Add routes in `client/src/App.tsx`
4. Update API client in `client/src/lib/api.ts`

### Styling

The app uses TailwindCSS with a custom design system defined in `client/src/index.css`. Color variables use OKLCH format for better color consistency.

### Type Safety

All API responses are typed in `client/src/lib/api.ts`. Add new types as needed for new features.

## License

MIT

## Support

For issues or questions, please contact the uduXPass development team.

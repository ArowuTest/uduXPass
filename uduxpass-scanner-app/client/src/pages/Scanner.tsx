/**
 * Design Philosophy: Professional Event Tech - Full-screen camera, animated scanning frame
 * Layout: Camera viewfinder with glowing blue frame, instructions card, manual entry fallback
 * Interaction: Auto-scan with haptic feedback
 */

import { useState, useEffect, useRef } from 'react';
import { useLocation } from 'wouter';
import { Button } from '@/components/ui/button';
import { Html5Qrcode } from 'html5-qrcode';
import { scannerApi, ScanningSession } from '@/lib/api';
import { ArrowLeft, Info } from 'lucide-react';
import { toast } from 'sonner';

export default function Scanner() {
  const [, setLocation] = useLocation();
  const [activeSession, setActiveSession] = useState<ScanningSession | null>(null);
  const [isScanning, setIsScanning] = useState(false);
  const [cameraError, setCameraError] = useState<string | null>(null);
  const html5QrCodeRef = useRef<Html5Qrcode | null>(null);
  const scannerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadActiveSession();
    return () => {
      stopScanner();
    };
  }, []);

  const loadActiveSession = async () => {
    try {
      const sessions = await scannerApi.getActiveSessions();
      if (sessions.length > 0) {
        setActiveSession(sessions[0]);
        startScanner();
      } else {
        toast.error('No active session found');
        setLocation('/dashboard');
      }
    } catch (error) {
      console.error('Failed to load active session:', error);
      toast.error('Failed to load session');
      setLocation('/dashboard');
    }
  };

  const startScanner = async () => {
    try {
      const html5QrCode = new Html5Qrcode('qr-reader');
      html5QrCodeRef.current = html5QrCode;

      await html5QrCode.start(
        { facingMode: 'environment' },
        {
          fps: 10,
          qrbox: { width: 250, height: 250 },
        },
        onScanSuccess,
        onScanError
      );

      setIsScanning(true);
      setCameraError(null);
    } catch (error: any) {
      console.error('Failed to start scanner:', error);
      setCameraError(error.message || 'Failed to access camera');
      toast.error('Camera access denied. Please enable camera permissions.');
    }
  };

  const stopScanner = async () => {
    if (html5QrCodeRef.current && isScanning) {
      try {
        await html5QrCodeRef.current.stop();
        html5QrCodeRef.current.clear();
      } catch (error) {
        console.error('Failed to stop scanner:', error);
      }
    }
  };

  const onScanSuccess = async (decodedText: string) => {
    if (!activeSession) return;

    // Stop scanner temporarily to prevent multiple scans
    await stopScanner();

    // Haptic feedback (if supported)
    if ('vibrate' in navigator) {
      navigator.vibrate(100);
    }

    try {
      const result = await scannerApi.validateTicket({
        qr_code_data: decodedText,
        session_id: activeSession.id,
      });

      // Backend returns: { success: boolean, message: string, data: { ticket: {...}, validated_at: ... } }
      if (result.success) {
        setLocation('/validation-success', { 
          state: { 
            ticket: result.data?.ticket,
            message: result.message || 'Ticket validated successfully'
          } 
        });
      } else {
        setLocation('/validation-error', { 
          state: { 
            message: result.error || result.message || 'Invalid ticket',
            ticket: result.data?.ticket 
          } 
        });
      }
    } catch (error: any) {
      console.error('Validation error:', error);
      const errorMessage = error.response?.data?.error || error.message || 'Failed to validate ticket';
      setLocation('/validation-error', { 
        state: { 
          message: errorMessage
        } 
      });
      // Restart scanner on error
      setTimeout(() => startScanner(), 1000);
    }
  };

  const onScanError = (errorMessage: string) => {
    // Ignore common scanning errors (no QR code in frame)
    if (!errorMessage.includes('NotFoundException')) {
      console.warn('Scan error:', errorMessage);
    }
  };

  const handleBack = async () => {
    await stopScanner();
    setLocation('/dashboard');
  };

  const handleManualEntry = async () => {
    await stopScanner();
    toast.info('Manual entry coming soon');
  };

  return (
    <div className="min-h-screen bg-black flex flex-col">
      {/* Top Bar */}
      <header className="bg-primary text-primary-foreground px-4 py-3 flex items-center justify-between z-10">
        <Button
          variant="ghost"
          size="sm"
          onClick={handleBack}
          className="text-primary-foreground hover:bg-primary-foreground/10"
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <h1 className="font-semibold text-lg">Scan Ticket</h1>
        <Button
          variant="ghost"
          size="sm"
          className="text-primary-foreground hover:bg-primary-foreground/10"
          onClick={() => toast.info('Position QR code within the frame. Camera will scan automatically.')}
        >
          <Info className="h-5 w-5" />
        </Button>
      </header>

      {/* Scanner Container */}
      <div className="flex-1 relative flex items-center justify-center">
        {/* QR Reader */}
        <div
          id="qr-reader"
          ref={scannerRef}
          className="w-full h-full"
          style={{ position: 'relative' }}
        />

        {/* Scanning Frame Overlay */}
        {isScanning && (
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="relative w-64 h-64">
              {/* Animated corners */}
              <div className="absolute top-0 left-0 w-12 h-12 border-t-4 border-l-4 border-primary rounded-tl-lg animate-pulse" />
              <div className="absolute top-0 right-0 w-12 h-12 border-t-4 border-r-4 border-primary rounded-tr-lg animate-pulse" />
              <div className="absolute bottom-0 left-0 w-12 h-12 border-b-4 border-l-4 border-primary rounded-bl-lg animate-pulse" />
              <div className="absolute bottom-0 right-0 w-12 h-12 border-b-4 border-r-4 border-primary rounded-br-lg animate-pulse" />
            </div>
          </div>
        )}

        {/* Camera Error */}
        {cameraError && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/90 p-6">
            <div className="text-center text-white">
              <p className="text-lg font-semibold mb-2">Camera Access Required</p>
              <p className="text-sm text-gray-300 mb-4">{cameraError}</p>
              <Button onClick={startScanner} variant="outline" className="text-white border-white hover:bg-white/10">
                Try Again
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Instructions Card */}
      <div className="bg-white rounded-t-3xl p-6 space-y-4">
        <div className="text-center">
          <p className="text-primary font-semibold mb-1">Position QR code within the frame</p>
          <p className="text-sm text-muted-foreground">Camera will scan automatically</p>
        </div>

        <Button
          variant="link"
          className="w-full text-primary"
          onClick={handleManualEntry}
        >
          Enter Ticket Code Manually
        </Button>
      </div>
    </div>
  );
}

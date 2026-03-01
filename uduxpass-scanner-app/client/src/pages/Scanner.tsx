/**
 * Scanner.tsx — uduXPass Scanner PWA
 * Design: Enterprise Amber/Slate — "Professional Event Operations"
 * - Full-screen camera with amber scan frame
 * - Processing overlay with spinner
 * - Manual entry bottom sheet
 * - Offline mode support
 * - Full API integration via ValidationResultContext
 */

import { useState, useEffect, useRef, useCallback } from 'react';
import { useLocation } from 'wouter';
import { Html5Qrcode } from 'html5-qrcode';
import { scannerApi, ScanningSession } from '@/lib/api';
import { ArrowLeft, WifiOff, Keyboard, Camera, Loader2, AlertCircle, X } from 'lucide-react';
import { toast } from 'sonner';
import { isOnline } from '@/lib/registerServiceWorker';
import { validateTicketOffline, shouldUseOfflineValidation } from '@/lib/offlineValidation';
import { useValidationResult } from '@/contexts/ValidationResultContext';

export default function Scanner() {
  const [, setLocation] = useLocation();
  const { setResult } = useValidationResult();

  // Refs to avoid stale closures
  const setResultRef = useRef(setResult);
  const setLocationRef = useRef(setLocation);
  useEffect(() => { setResultRef.current = setResult; }, [setResult]);
  useEffect(() => { setLocationRef.current = setLocation; }, [setLocation]);

  const [activeSession, setActiveSession] = useState<ScanningSession | null>(null);
  const [isScanning, setIsScanning] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [cameraError, setCameraError] = useState<string | null>(null);
  const [offline, setOffline] = useState(!isOnline());
  const [showManual, setShowManual] = useState(false);
  const [manualCode, setManualCode] = useState('');
  const [isSubmittingManual, setIsSubmittingManual] = useState(false);
  const html5QrCodeRef = useRef<Html5Qrcode | null>(null);
  const scannerDivId = 'qr-reader-canvas';
  const processingRef = useRef(false);

  useEffect(() => {
    const on = () => { setOffline(false); toast.success('Back online — syncing data'); };
    const off = () => { setOffline(true); toast.warning('Offline mode — validations will sync when online'); };
    window.addEventListener('online', on);
    window.addEventListener('offline', off);
    return () => { window.removeEventListener('online', on); window.removeEventListener('offline', off); };
  }, []);

  useEffect(() => {
    loadActiveSession();
    return () => { stopScanner(); };
  }, []);

  const loadActiveSession = async () => {
    try {
      const session = await scannerApi.getCurrentSession();
      if (session) {
        setActiveSession(session);
        setTimeout(() => startScanner(), 300);
      } else {
        toast.error('No active session. Please start a session first.');
        setLocation('/dashboard');
      }
    } catch {
      toast.error('Failed to load session');
      setLocation('/dashboard');
    }
  };

  const startScanner = async () => {
    try {
      if (html5QrCodeRef.current) {
        try { await html5QrCodeRef.current.stop(); } catch {}
        html5QrCodeRef.current = null;
      }
      const qr = new Html5Qrcode(scannerDivId);
      html5QrCodeRef.current = qr;
      await qr.start(
        { facingMode: 'environment' },
        { fps: 10, qrbox: { width: 240, height: 240 } },
        onScanSuccess,
        onScanError
      );
      setIsScanning(true);
      setCameraError(null);
    } catch (err: any) {
      setCameraError(err?.message || 'Failed to access camera');
    }
  };

  const stopScanner = async () => {
    if (html5QrCodeRef.current) {
      try { await html5QrCodeRef.current.stop(); html5QrCodeRef.current.clear(); } catch {}
      html5QrCodeRef.current = null;
    }
    setIsScanning(false);
  };

  const processValidation = useCallback(async (decodedText: string) => {
    if (!activeSession || processingRef.current) return;
    processingRef.current = true;
    setIsProcessing(true);
    if ('vibrate' in navigator) navigator.vibrate(100);

    try {
      let result: any;
      if (shouldUseOfflineValidation()) {
        result = await validateTicketOffline(decodedText, activeSession.id);
        if (result.offline) toast.info('Offline mode — validation will sync when online');
      } else {
        result = await scannerApi.validateTicket({
          ticket_code: decodedText,
          event_id: activeSession.event_id,
        });
      }

      if (result.success && result.valid) {
        setResultRef.current({
          success: true,
          valid: true,
          already_validated: result.already_validated || false,
          message: result.message || 'Ticket validated successfully',
          serial_number: result.serial_number,
          ticket_id: result.ticket_id,
          holder_name: result.holder_name,
          ticket_tier: result.ticket_type,
          validated_at: result.validation_time,
        });
        setLocationRef.current('/validation-success');
      } else {
        setResultRef.current({
          success: false,
          valid: false,
          already_validated: result.already_validated || false,
          message: result.message || 'Invalid ticket',
          serial_number: result.serial_number,
          ticket_id: result.ticket_id,
          error: result.already_validated ? 'ALREADY_VALIDATED' : 'INVALID',
        });
        setLocationRef.current('/validation-error');
      }
    } catch (err: any) {
      const msg = err?.response?.data?.error || err?.message || 'Failed to validate ticket';
      setResultRef.current({
        success: false,
        valid: false,
        already_validated: false,
        message: msg,
        error: 'SYSTEM_ERROR',
      });
      setLocationRef.current('/validation-error');
    } finally {
      processingRef.current = false;
      setIsProcessing(false);
    }
  }, [activeSession]);

  const onScanSuccess = async (decodedText: string) => {
    await stopScanner();
    await processValidation(decodedText);
  };

  const onScanError = (msg: string) => {
    if (!msg.includes('NotFoundException')) console.warn('Scan error:', msg);
  };

  const handleManualSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!manualCode.trim()) return;
    setIsSubmittingManual(true);
    await stopScanner();
    await processValidation(manualCode.trim());
    setManualCode('');
    setIsSubmittingManual(false);
  };

  const handleBack = async () => {
    await stopScanner();
    setLocation('/dashboard');
  };

  return (
    <div className="min-h-screen bg-black flex flex-col relative overflow-hidden">

      {/* Floating Header */}
      <header
        className="absolute top-0 left-0 right-0 z-30 flex items-center justify-between px-4 py-4"
        style={{
          background: 'linear-gradient(to bottom, rgba(0,0,0,0.75) 0%, transparent 100%)',
        }}
      >
        <button
          onClick={handleBack}
          className="w-10 h-10 rounded-2xl flex items-center justify-center"
          style={{
            background: 'rgba(255,255,255,0.12)',
            backdropFilter: 'blur(8px)',
            border: '1px solid rgba(255,255,255,0.15)',
            color: 'white',
          }}
        >
          <ArrowLeft size={18} />
        </button>

        <div className="text-center">
          <p
            style={{
              fontFamily: 'Space Grotesk, sans-serif',
              fontSize: '15px',
              fontWeight: 700,
              color: 'white',
              lineHeight: 1.2,
            }}
          >
            Scan Ticket
          </p>
          {activeSession?.event_name && (
            <p style={{ fontSize: '11px', color: 'rgba(255,255,255,0.65)', lineHeight: 1.2 }}>
              {activeSession.event_name}
            </p>
          )}
        </div>

        <button
          onClick={() => setShowManual(!showManual)}
          className="w-10 h-10 rounded-2xl flex items-center justify-center"
          style={{
            background: showManual ? 'rgba(245,158,11,0.3)' : 'rgba(255,255,255,0.12)',
            backdropFilter: 'blur(8px)',
            border: `1px solid ${showManual ? 'rgba(245,158,11,0.5)' : 'rgba(255,255,255,0.15)'}`,
            color: showManual ? '#F59E0B' : 'white',
          }}
        >
          <Keyboard size={17} />
        </button>
      </header>

      {/* Camera / QR Reader */}
      <div className="flex-1 relative">
        <div
          id={scannerDivId}
          className="w-full"
          style={{ minHeight: '100vh' }}
        />

        {/* Scan frame overlay */}
        {isScanning && !isProcessing && !showManual && (
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            {/* Vignette */}
            <div
              className="absolute inset-0"
              style={{
                background: 'radial-gradient(ellipse 280px 280px at center, transparent 0%, rgba(0,0,0,0.65) 100%)',
              }}
            />
            {/* Amber corner brackets */}
            <div className="relative w-64 h-64">
              {/* TL */}
              <div className="absolute top-0 left-0 w-12 h-12" style={{ borderTop: '3px solid #F59E0B', borderLeft: '3px solid #F59E0B', borderRadius: '6px 0 0 0' }} />
              {/* TR */}
              <div className="absolute top-0 right-0 w-12 h-12" style={{ borderTop: '3px solid #F59E0B', borderRight: '3px solid #F59E0B', borderRadius: '0 6px 0 0' }} />
              {/* BL */}
              <div className="absolute bottom-0 left-0 w-12 h-12" style={{ borderBottom: '3px solid #F59E0B', borderLeft: '3px solid #F59E0B', borderRadius: '0 0 0 6px' }} />
              {/* BR */}
              <div className="absolute bottom-0 right-0 w-12 h-12" style={{ borderBottom: '3px solid #F59E0B', borderRight: '3px solid #F59E0B', borderRadius: '0 0 6px 0' }} />
              {/* Scan line */}
              <div
                className="absolute left-3 right-3 h-0.5"
                style={{
                  background: 'linear-gradient(90deg, transparent, #F59E0B, transparent)',
                  animation: 'scanLine 2s ease-in-out infinite',
                }}
              />
            </div>
          </div>
        )}

        {/* Processing overlay */}
        {isProcessing && (
          <div className="absolute inset-0 flex items-center justify-center z-20" style={{ background: 'rgba(0,0,0,0.85)' }}>
            <div className="text-center">
              <div
                className="w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-4"
                style={{
                  background: 'rgba(245,158,11,0.15)',
                  border: '2px solid rgba(245,158,11,0.4)',
                  boxShadow: '0 0 40px rgba(245,158,11,0.2)',
                }}
              >
                <Loader2 size={36} style={{ color: '#F59E0B', animation: 'spin 0.8s linear infinite' }} />
              </div>
              <p
                style={{
                  fontFamily: 'Space Grotesk, sans-serif',
                  fontSize: '16px',
                  fontWeight: 700,
                  color: 'white',
                  marginBottom: '4px',
                }}
              >
                Validating...
              </p>
              <p style={{ fontSize: '13px', color: 'rgba(255,255,255,0.5)' }}>Checking ticket authenticity</p>
            </div>
          </div>
        )}

        {/* Camera error overlay — hide when manual is active */}
        {cameraError && !showManual && (
          <div className="absolute inset-0 flex items-center justify-center z-20 p-6" style={{ background: 'rgba(0,0,0,0.92)' }}>
            <div className="text-center max-w-xs w-full">
              <div
                className="w-16 h-16 rounded-3xl flex items-center justify-center mx-auto mb-4"
                style={{ background: 'rgba(244,63,94,0.15)', border: '1px solid rgba(244,63,94,0.3)' }}
              >
                <AlertCircle size={28} style={{ color: '#F43F5E' }} />
              </div>
              <p
                style={{
                  fontFamily: 'Space Grotesk, sans-serif',
                  fontSize: '18px',
                  fontWeight: 700,
                  color: 'white',
                  marginBottom: '8px',
                }}
              >
                Camera Access Required
              </p>
              <p style={{ fontSize: '13px', color: 'rgba(255,255,255,0.55)', marginBottom: '24px', lineHeight: 1.5 }}>
                {cameraError}
              </p>
              <div className="space-y-2">
                <button
                  onClick={startScanner}
                  className="w-full h-12 rounded-2xl flex items-center justify-center gap-2 font-semibold"
                  style={{
                    background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)',
                    color: 'oklch(0.13 0.025 245)',
                    fontSize: '14px',
                    fontFamily: 'Space Grotesk, sans-serif',
                  }}
                >
                  <Camera size={18} />
                  Try Again
                </button>
                <button
                  onClick={() => setShowManual(true)}
                  className="w-full h-12 rounded-2xl flex items-center justify-center gap-2 font-semibold"
                  style={{
                    background: 'rgba(255,255,255,0.08)',
                    border: '1px solid rgba(255,255,255,0.15)',
                    color: 'rgba(255,255,255,0.85)',
                    fontSize: '14px',
                    fontFamily: 'Space Grotesk, sans-serif',
                  }}
                >
                  <Keyboard size={18} />
                  Enter Code Manually
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Offline badge */}
        {offline && (
          <div
            className="absolute top-20 left-1/2 -translate-x-1/2 flex items-center gap-1.5 px-3 py-1.5 rounded-full z-20"
            style={{
              background: 'rgba(245,158,11,0.2)',
              border: '1px solid rgba(245,158,11,0.4)',
              backdropFilter: 'blur(8px)',
            }}
          >
            <WifiOff size={12} style={{ color: '#F59E0B' }} />
            <span style={{ fontSize: '11px', fontWeight: 600, color: '#F59E0B' }}>Offline Mode</span>
          </div>
        )}
      </div>

      {/* Bottom sheet — instructions or manual entry */}
      <div
        className="relative z-20"
        style={{
          background: 'oklch(0.16 0.025 245)',
          borderTop: '1px solid oklch(0.24 0.025 245)',
          borderRadius: '24px 24px 0 0',
          padding: showManual ? '20px 20px 32px' : '16px 20px 32px',
        }}
      >
        {/* Drag handle */}
        <div
          className="w-10 h-1 rounded-full mx-auto mb-4"
          style={{ background: 'oklch(0.30 0.025 245)' }}
        />

        {showManual ? (
          <form onSubmit={handleManualSubmit}>
            <div className="flex items-center justify-between mb-3">
              <p
                style={{
                  fontFamily: 'Space Grotesk, sans-serif',
                  fontSize: '15px',
                  fontWeight: 700,
                  color: 'oklch(0.97 0.005 240)',
                }}
              >
                Enter Ticket Code
              </p>
              <button
                type="button"
                onClick={() => { setShowManual(false); setManualCode(''); }}
                className="w-7 h-7 rounded-lg flex items-center justify-center"
                style={{ background: 'oklch(0.24 0.025 245)', color: 'oklch(0.55 0.015 240)' }}
              >
                <X size={14} />
              </button>
            </div>

            <textarea
              value={manualCode}
              onChange={(e) => setManualCode(e.target.value)}
              placeholder="Paste or type the ticket QR code (JWT)..."
              autoFocus
              rows={3}
              className="w-full rounded-2xl px-4 py-3 font-mono resize-none mb-3"
              style={{
                background: 'oklch(0.20 0.025 245)',
                border: '1px solid oklch(0.28 0.025 245)',
                color: 'oklch(0.97 0.005 240)',
                fontSize: '12px',
                outline: 'none',
                lineHeight: 1.5,
              }}
            />

            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => { setShowManual(false); setManualCode(''); }}
                className="flex-1 h-12 rounded-2xl font-semibold"
                style={{
                  background: 'oklch(0.22 0.025 245)',
                  border: '1px solid oklch(0.28 0.025 245)',
                  color: 'oklch(0.70 0.015 240)',
                  fontSize: '14px',
                  fontFamily: 'Space Grotesk, sans-serif',
                }}
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={!manualCode.trim() || isSubmittingManual}
                className="flex-1 h-12 rounded-2xl flex items-center justify-center gap-2 font-semibold"
                style={{
                  background: manualCode.trim() ? 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)' : 'oklch(0.22 0.025 245)',
                  color: manualCode.trim() ? 'oklch(0.13 0.025 245)' : 'oklch(0.45 0.015 240)',
                  fontSize: '14px',
                  fontFamily: 'Space Grotesk, sans-serif',
                  transition: 'all 150ms ease',
                }}
              >
                {isSubmittingManual ? (
                  <Loader2 size={18} className="animate-spin" />
                ) : (
                  'Validate'
                )}
              </button>
            </div>
          </form>
        ) : (
          <div className="text-center">
            <p
              style={{
                fontFamily: 'Space Grotesk, sans-serif',
                fontSize: '14px',
                fontWeight: 600,
                color: 'oklch(0.97 0.005 240)',
                marginBottom: '4px',
              }}
            >
              {isScanning ? 'Position QR code within the frame' : 'Initialising camera...'}
            </p>
            <p style={{ fontSize: '12px', color: 'oklch(0.50 0.015 240)', marginBottom: '12px' }}>
              Camera will scan automatically
            </p>
            <button
              onClick={() => setShowManual(true)}
              style={{
                fontSize: '12px',
                fontWeight: 600,
                color: '#F59E0B',
                background: 'none',
                border: 'none',
                cursor: 'pointer',
              }}
            >
              Enter code manually
            </button>
          </div>
        )}
      </div>

      <style>{`
        @keyframes scanLine {
          0% { top: 10%; opacity: 0; }
          10% { opacity: 1; }
          90% { opacity: 1; }
          100% { top: 90%; opacity: 0; }
        }
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}

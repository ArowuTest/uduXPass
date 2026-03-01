/**
 * ValidationError.tsx — uduXPass Scanner PWA
 * Design: Enterprise — Full-screen rose red (INVALID) or amber (DUPLICATE)
 * - Unmistakable DENIED state readable from arm's length
 * - Error details card with reason
 * - Scan Next / Dashboard CTAs
 */

import { useEffect, useRef } from 'react';
import { useLocation } from 'wouter';
import { useValidationResult } from '@/contexts/ValidationResultContext';
import { XCircle, AlertTriangle, Scan, LayoutDashboard, Hash, AlertCircle } from 'lucide-react';

export default function ValidationError() {
  const [, setLocation] = useLocation();
  const { result, clearResult } = useValidationResult();
  const resultRef = useRef(result);

  const data = resultRef.current;
  const isAlreadyValidated = data?.already_validated || data?.error === 'ALREADY_VALIDATED';
  const isSystemError = data?.error === 'SYSTEM_ERROR';

  useEffect(() => {
    if ('vibrate' in navigator) navigator.vibrate([200, 100, 200]);
  }, []);

  const handleScanNext = () => {
    clearResult();
    setLocation('/scan');
  };

  const handleDashboard = () => {
    clearResult();
    setLocation('/dashboard');
  };

  // Color config based on error type
  const cfg = isAlreadyValidated
    ? {
        bg: 'oklch(0.17 0.08 60)',
        glow: 'rgba(245,158,11,0.28)',
        accent: '#F59E0B',
        accentBg: 'rgba(245,158,11,0.18)',
        accentBorder: 'rgba(245,158,11,0.40)',
        badge: 'DUPLICATE SCAN',
        heading: 'Already Scanned',
        icon: <AlertTriangle size={52} style={{ color: '#F59E0B' }} strokeWidth={1.5} />,
      }
    : isSystemError
    ? {
        bg: 'oklch(0.16 0.03 240)',
        glow: 'rgba(148,163,184,0.20)',
        accent: '#94A3B8',
        accentBg: 'rgba(148,163,184,0.12)',
        accentBorder: 'rgba(148,163,184,0.30)',
        badge: 'SYSTEM ERROR',
        heading: 'System Error',
        icon: <AlertCircle size={52} style={{ color: '#94A3B8' }} strokeWidth={1.5} />,
      }
    : {
        bg: 'oklch(0.15 0.08 15)',
        glow: 'rgba(244,63,94,0.28)',
        accent: '#F43F5E',
        accentBg: 'rgba(244,63,94,0.18)',
        accentBorder: 'rgba(244,63,94,0.40)',
        badge: 'INVALID',
        heading: 'Invalid Ticket',
        icon: <XCircle size={52} style={{ color: '#F43F5E' }} strokeWidth={1.5} />,
      };

  return (
    <div
      className="min-h-screen flex flex-col relative overflow-hidden"
      style={{ background: cfg.bg }}
    >
      {/* Radial glow */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: `radial-gradient(ellipse 80% 55% at 50% 25%, ${cfg.glow} 0%, transparent 70%)`,
        }}
      />

      <div className="relative flex-1 flex flex-col items-center justify-center px-6 py-10 text-center">

        {/* Icon */}
        <div className="relative mb-6">
          <div
            className="absolute inset-0 rounded-full"
            style={{
              background: cfg.accentBg,
              animation: 'pulse 2s ease-in-out infinite',
            }}
          />
          <div
            className="relative w-28 h-28 rounded-full flex items-center justify-center"
            style={{
              background: cfg.accentBg,
              border: `2.5px solid ${cfg.accentBorder}`,
              boxShadow: `0 0 60px ${cfg.glow}`,
            }}
          >
            {cfg.icon}
          </div>
        </div>

        {/* Badge */}
        <div
          className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full mb-4"
          style={{
            background: cfg.accentBg,
            border: `1px solid ${cfg.accentBorder}`,
          }}
        >
          <div
            className="w-2 h-2 rounded-full"
            style={{ background: cfg.accent, boxShadow: `0 0 6px ${cfg.accent}` }}
          />
          <span style={{ fontSize: '11px', fontWeight: 700, color: cfg.accent, letterSpacing: '0.12em', textTransform: 'uppercase' }}>
            {cfg.badge}
          </span>
        </div>

        <h1
          style={{
            fontFamily: 'Space Grotesk, sans-serif',
            fontSize: '2.75rem',
            fontWeight: 900,
            color: 'white',
            letterSpacing: '-0.04em',
            lineHeight: 1,
            marginBottom: '8px',
          }}
        >
          {isAlreadyValidated ? 'DUPLICATE' : isSystemError ? 'ERROR' : 'DENIED'}
        </h1>
        <p style={{ fontSize: '14px', color: 'rgba(255,255,255,0.60)', marginBottom: '28px', maxWidth: '280px' }}>
          {data?.message || 'This ticket could not be validated'}
        </p>

        {/* Error details card */}
        <div
          className="w-full max-w-sm rounded-3xl overflow-hidden mb-8"
          style={{
            background: 'rgba(0,0,0,0.22)',
            border: `1px solid ${cfg.accentBorder}`,
            backdropFilter: 'blur(12px)',
          }}
        >
          <div
            className="px-5 py-3 flex items-center gap-2"
            style={{ borderBottom: `1px solid ${cfg.accentBorder}` }}
          >
            <div className="w-1.5 h-1.5 rounded-full" style={{ background: cfg.accent }} />
            <span style={{ fontSize: '11px', fontWeight: 700, color: 'rgba(255,255,255,0.45)', letterSpacing: '0.1em', textTransform: 'uppercase' }}>
              Error Details
            </span>
          </div>

          <div className="px-5 py-4 space-y-3">
            <div>
              <p style={{ fontSize: '10px', color: 'rgba(255,255,255,0.38)', textTransform: 'uppercase', letterSpacing: '0.08em', marginBottom: '4px' }}>
                Reason
              </p>
              <p style={{ fontSize: '14px', fontWeight: 600, color: cfg.accent }}>
                {data?.message || 'Invalid ticket'}
              </p>
            </div>

            {data?.serial_number && (
              <div
                className="flex items-center gap-3 pt-3"
                style={{ borderTop: '1px solid rgba(255,255,255,0.08)' }}
              >
                <Hash size={14} style={{ color: cfg.accent, flexShrink: 0 }} />
                <div className="text-left">
                  <p style={{ fontSize: '10px', color: 'rgba(255,255,255,0.38)', textTransform: 'uppercase', letterSpacing: '0.08em' }}>Serial Number</p>
                  <p style={{ fontFamily: 'monospace', fontSize: '14px', fontWeight: 700, color: 'white' }}>{data.serial_number}</p>
                </div>
              </div>
            )}

            {isAlreadyValidated && (
              <div
                className="pt-3"
                style={{ borderTop: '1px solid rgba(255,255,255,0.08)' }}
              >
                <p style={{ fontSize: '12px', color: 'rgba(255,255,255,0.45)', lineHeight: 1.5 }}>
                  This ticket has already been used for entry. Each ticket is valid for a single use only.
                </p>
              </div>
            )}
          </div>
        </div>

        {/* CTAs */}
        <div className="w-full max-w-sm space-y-2">
          <button
            onClick={handleScanNext}
            className="w-full flex items-center justify-center gap-2 font-semibold"
            style={{
              height: '52px',
              borderRadius: '16px',
              background: cfg.accent,
              color: isAlreadyValidated ? 'oklch(0.13 0.025 245)' : 'white',
              fontSize: '15px',
              fontFamily: 'Space Grotesk, sans-serif',
              boxShadow: `0 4px 20px ${cfg.glow}`,
              border: 'none',
              cursor: 'pointer',
            }}
          >
            <Scan size={20} />
            Scan Next Ticket
          </button>
          <button
            onClick={handleDashboard}
            className="w-full flex items-center justify-center gap-2"
            style={{
              height: '48px',
              borderRadius: '16px',
              background: 'rgba(255,255,255,0.07)',
              border: '1px solid rgba(255,255,255,0.14)',
              color: 'rgba(255,255,255,0.70)',
              fontSize: '14px',
              fontFamily: 'Space Grotesk, sans-serif',
              cursor: 'pointer',
            }}
          >
            <LayoutDashboard size={17} />
            Back to Dashboard
          </button>
        </div>
      </div>

      <style>{`
        @keyframes pulse {
          0%, 100% { transform: scale(1); opacity: 0.6; }
          50% { transform: scale(1.15); opacity: 0.2; }
        }
      `}</style>
    </div>
  );
}

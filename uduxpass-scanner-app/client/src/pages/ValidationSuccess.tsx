/**
 * ValidationSuccess.tsx — uduXPass Scanner PWA
 * Design: Enterprise — Full-screen emerald green ADMITTED state
 * - Unmistakable VALID state readable from arm's length
 * - Ticket details card with real API data
 * - Auto-redirect to scanner after 5s
 */

import { useEffect, useRef, useState } from 'react';
import { useLocation } from 'wouter';
import { useValidationResult } from '@/contexts/ValidationResultContext';
import { CheckCircle2, Scan, LayoutDashboard, Clock, Hash, User, Tag } from 'lucide-react';

export default function ValidationSuccess() {
  const [, setLocation] = useLocation();
  const { result, clearResult } = useValidationResult();
  const [countdown, setCountdown] = useState(5);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const redirectedRef = useRef(false);
  const resultRef = useRef(result);

  useEffect(() => {
    if (!resultRef.current) {
      setLocation('/scan');
      return;
    }
    if ('vibrate' in navigator) navigator.vibrate([80, 40, 80]);

    timerRef.current = setInterval(() => {
      setCountdown((c) => {
        if (c <= 1) {
          if (!redirectedRef.current) {
            redirectedRef.current = true;
            clearResult();
            setLocation('/scan');
          }
          return 0;
        }
        return c - 1;
      });
    }, 1000);

    return () => { if (timerRef.current) clearInterval(timerRef.current); };
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleScanNext = () => {
    if (timerRef.current) clearInterval(timerRef.current);
    clearResult();
    setLocation('/scan');
  };

  const handleDashboard = () => {
    if (timerRef.current) clearInterval(timerRef.current);
    clearResult();
    setLocation('/dashboard');
  };

  const formatTime = (iso?: string) => {
    if (!iso) return null;
    try {
      return new Date(iso).toLocaleTimeString('en-NG', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    } catch { return iso; }
  };

  const data = resultRef.current;

  return (
    <div
      className="min-h-screen flex flex-col relative overflow-hidden"
      style={{ background: 'oklch(0.15 0.10 155)' }}
    >
      {/* Radial glow */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: 'radial-gradient(ellipse 80% 55% at 50% 25%, rgba(16,185,129,0.30) 0%, transparent 70%)',
        }}
      />

      <div className="relative flex-1 flex flex-col items-center justify-center px-6 py-10 text-center">

        {/* Big check icon */}
        <div className="relative mb-6">
          <div
            className="absolute inset-0 rounded-full"
            style={{
              background: 'rgba(16,185,129,0.15)',
              animation: 'pulse 2s ease-in-out infinite',
            }}
          />
          <div
            className="relative w-28 h-28 rounded-full flex items-center justify-center"
            style={{
              background: 'rgba(16,185,129,0.18)',
              border: '2.5px solid rgba(16,185,129,0.55)',
              boxShadow: '0 0 60px rgba(16,185,129,0.35)',
            }}
          >
            <CheckCircle2 size={56} style={{ color: '#10B981' }} strokeWidth={1.5} />
          </div>
        </div>

        {/* VALID badge */}
        <div
          className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full mb-4"
          style={{
            background: 'rgba(16,185,129,0.18)',
            border: '1px solid rgba(16,185,129,0.40)',
          }}
        >
          <div
            className="w-2 h-2 rounded-full"
            style={{ background: '#10B981', boxShadow: '0 0 6px rgba(16,185,129,0.8)' }}
          />
          <span style={{ fontSize: '11px', fontWeight: 700, color: '#10B981', letterSpacing: '0.12em', textTransform: 'uppercase' }}>
            Valid Ticket
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
          ADMITTED
        </h1>
        <p style={{ fontSize: '14px', color: 'rgba(255,255,255,0.60)', marginBottom: '28px' }}>
          {data?.message || 'Ticket validated successfully'}
        </p>

        {/* Ticket details card */}
        {data && (
          <div
            className="w-full max-w-sm rounded-3xl overflow-hidden mb-8"
            style={{
              background: 'rgba(0,0,0,0.22)',
              border: '1px solid rgba(16,185,129,0.22)',
              backdropFilter: 'blur(12px)',
            }}
          >
            <div
              className="px-5 py-3 flex items-center gap-2"
              style={{ borderBottom: '1px solid rgba(16,185,129,0.15)' }}
            >
              <div className="w-1.5 h-1.5 rounded-full" style={{ background: '#10B981' }} />
              <span style={{ fontSize: '11px', fontWeight: 700, color: 'rgba(255,255,255,0.45)', letterSpacing: '0.1em', textTransform: 'uppercase' }}>
                Ticket Details
              </span>
            </div>

            <div className="px-5 py-4 space-y-3">
              {data.serial_number && (
                <div className="flex items-center gap-3">
                  <Hash size={14} style={{ color: '#10B981', flexShrink: 0 }} />
                  <div className="text-left">
                    <p style={{ fontSize: '10px', color: 'rgba(255,255,255,0.38)', textTransform: 'uppercase', letterSpacing: '0.08em' }}>Serial Number</p>
                    <p style={{ fontFamily: 'monospace', fontSize: '14px', fontWeight: 700, color: 'white' }}>{data.serial_number}</p>
                  </div>
                </div>
              )}
              {data.holder_name && (
                <div className="flex items-center gap-3">
                  <User size={14} style={{ color: '#10B981', flexShrink: 0 }} />
                  <div className="text-left">
                    <p style={{ fontSize: '10px', color: 'rgba(255,255,255,0.38)', textTransform: 'uppercase', letterSpacing: '0.08em' }}>Ticket Holder</p>
                    <p style={{ fontSize: '14px', fontWeight: 600, color: 'white' }}>{data.holder_name}</p>
                  </div>
                </div>
              )}
              {data.ticket_tier && (
                <div className="flex items-center gap-3">
                  <Tag size={14} style={{ color: '#10B981', flexShrink: 0 }} />
                  <div className="text-left">
                    <p style={{ fontSize: '10px', color: 'rgba(255,255,255,0.38)', textTransform: 'uppercase', letterSpacing: '0.08em' }}>Ticket Tier</p>
                    <p style={{ fontSize: '14px', fontWeight: 600, color: 'white' }}>{data.ticket_tier}</p>
                  </div>
                </div>
              )}
              {data.validated_at && (
                <div className="flex items-center gap-3">
                  <Clock size={14} style={{ color: '#10B981', flexShrink: 0 }} />
                  <div className="text-left">
                    <p style={{ fontSize: '10px', color: 'rgba(255,255,255,0.38)', textTransform: 'uppercase', letterSpacing: '0.08em' }}>Validated At</p>
                    <p style={{ fontSize: '14px', fontWeight: 600, color: 'white' }}>{formatTime(data.validated_at)}</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Countdown */}
        <p style={{ fontSize: '13px', color: 'rgba(255,255,255,0.40)', marginBottom: '20px' }}>
          Returning to scanner in{' '}
          <span style={{ color: '#10B981', fontWeight: 700 }}>{countdown}s</span>
        </p>

        {/* CTAs */}
        <div className="w-full max-w-sm space-y-2">
          <button
            onClick={handleScanNext}
            className="w-full flex items-center justify-center gap-2 font-semibold"
            style={{
              height: '52px',
              borderRadius: '16px',
              background: '#10B981',
              color: 'oklch(0.13 0.025 245)',
              fontSize: '15px',
              fontFamily: 'Space Grotesk, sans-serif',
              boxShadow: '0 4px 20px rgba(16,185,129,0.40)',
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

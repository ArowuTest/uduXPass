/**
 * SessionHistory.tsx — uduXPass Scanner PWA
 * Design: Enterprise Amber/Slate — "Professional Event Operations"
 * - Session cards with real stats from API
 * - Active session highlighted with amber accent
 * - Duration, scan counts, valid/invalid breakdown
 */

import { useState, useEffect } from 'react';
import { useLocation } from 'wouter';
import { scannerApi, ScanningSession } from '@/lib/api';
import {
  ArrowLeft,
  Scan,
  CheckCircle2,
  XCircle,
  Loader2,
  Clock,
  Activity,
  MapPin,
  CalendarDays,
  History,
} from 'lucide-react';
import { toast } from 'sonner';

export default function SessionHistory() {
  const [, setLocation] = useLocation();
  const [sessions, setSessions] = useState<ScanningSession[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadSessions();
  }, []);

  const loadSessions = async () => {
    try {
      const session = await scannerApi.getCurrentSession();
      setSessions(session ? [session] : []);
    } catch {
      toast.error('Failed to load session history');
    } finally {
      setIsLoading(false);
    }
  };

  const formatDateTime = (iso: string) => {
    try {
      return new Date(iso).toLocaleString('en-NG', {
        day: 'numeric', month: 'short', year: 'numeric',
        hour: '2-digit', minute: '2-digit',
      });
    } catch { return iso; }
  };

  const getDuration = (start: string, end?: string | null) => {
    try {
      const s = new Date(start).getTime();
      const e = end ? new Date(end).getTime() : Date.now();
      const diff = Math.floor((e - s) / 1000);
      const h = Math.floor(diff / 3600);
      const m = Math.floor((diff % 3600) / 60);
      if (h > 0) return `${h}h ${m}m`;
      return `${m}m`;
    } catch { return '—'; }
  };

  const getSuccessRate = (valid: number, total: number) => {
    if (!total) return 0;
    return Math.round((valid / total) * 100);
  };

  return (
    <div
      className="min-h-screen flex flex-col"
      style={{ background: 'oklch(0.13 0.025 245)' }}
    >
      {/* Sticky header */}
      <header
        className="sticky top-0 z-20 flex items-center gap-3 px-4 py-4"
        style={{
          background: 'oklch(0.16 0.025 245)',
          borderBottom: '1px solid oklch(0.22 0.025 245)',
          backdropFilter: 'blur(12px)',
        }}
      >
        <button
          onClick={() => setLocation('/dashboard')}
          className="w-9 h-9 rounded-xl flex items-center justify-center"
          style={{
            background: 'oklch(0.22 0.025 245)',
            color: 'oklch(0.65 0.015 240)',
            border: 'none',
            cursor: 'pointer',
          }}
        >
          <ArrowLeft size={17} />
        </button>
        <div>
          <h1
            style={{
              fontFamily: 'Space Grotesk, sans-serif',
              fontSize: '15px',
              fontWeight: 700,
              color: 'oklch(0.97 0.005 240)',
              lineHeight: 1.2,
            }}
          >
            Session History
          </h1>
          <p style={{ fontSize: '11px', color: 'oklch(0.50 0.015 240)', lineHeight: 1.2 }}>
            Your scanning sessions
          </p>
        </div>
        <div className="ml-auto">
          <History size={18} style={{ color: 'oklch(0.40 0.015 240)' }} />
        </div>
      </header>

      {/* Content */}
      <main className="flex-1 px-4 py-5 max-w-lg mx-auto w-full">
        {isLoading ? (
          <div className="flex items-center justify-center py-20">
            <div className="text-center">
              <Loader2 size={28} style={{ color: '#F59E0B', animation: 'spin 0.8s linear infinite', margin: '0 auto 12px' }} />
              <p style={{ fontSize: '13px', color: 'oklch(0.50 0.015 240)' }}>Loading sessions...</p>
            </div>
          </div>
        ) : sessions.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-center">
            <div
              className="w-16 h-16 rounded-3xl flex items-center justify-center mb-4"
              style={{ background: 'oklch(0.20 0.025 245)', border: '1px solid oklch(0.26 0.025 245)' }}
            >
              <Scan size={28} style={{ color: 'oklch(0.45 0.015 240)' }} />
            </div>
            <p
              style={{
                fontFamily: 'Space Grotesk, sans-serif',
                fontSize: '16px',
                fontWeight: 700,
                color: 'oklch(0.97 0.005 240)',
                marginBottom: '6px',
              }}
            >
              No sessions yet
            </p>
            <p style={{ fontSize: '13px', color: 'oklch(0.50 0.015 240)', marginBottom: '20px', maxWidth: '240px', lineHeight: 1.5 }}>
              Start a new scanning session to begin validating tickets.
            </p>
            <button
              onClick={() => setLocation('/create-session')}
              className="flex items-center gap-2 font-semibold px-5"
              style={{
                height: '44px',
                borderRadius: '14px',
                background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)',
                color: 'oklch(0.13 0.025 245)',
                fontSize: '14px',
                fontFamily: 'Space Grotesk, sans-serif',
                border: 'none',
                cursor: 'pointer',
              }}
            >
              <Scan size={16} />
              Start New Session
            </button>
          </div>
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => {
              const isActive = session.status === 'active' || !session.end_time;
              const total = session.scans_count ?? 0;
              const valid = session.valid_scans ?? 0;
              const invalid = session.invalid_scans ?? 0;
              const rate = getSuccessRate(valid, total);

              return (
                <div
                  key={session.id}
                  className="rounded-3xl overflow-hidden"
                  style={{
                    background: 'oklch(0.17 0.025 245)',
                    border: isActive
                      ? '1px solid rgba(245,158,11,0.35)'
                      : '1px solid oklch(0.24 0.025 245)',
                    boxShadow: isActive ? '0 4px 24px rgba(245,158,11,0.08)' : 'none',
                  }}
                >
                  {/* Card header */}
                  <div
                    className="px-5 py-3 flex items-center justify-between"
                    style={{
                      background: isActive ? 'rgba(245,158,11,0.08)' : 'oklch(0.20 0.025 245)',
                      borderBottom: isActive
                        ? '1px solid rgba(245,158,11,0.20)'
                        : '1px solid oklch(0.24 0.025 245)',
                    }}
                  >
                    <div className="flex items-center gap-2">
                      {isActive && (
                        <div
                          className="w-2 h-2 rounded-full"
                          style={{ background: '#F59E0B', boxShadow: '0 0 6px rgba(245,158,11,0.6)', animation: 'pulse 2s infinite' }}
                        />
                      )}
                      <span
                        style={{
                          fontSize: '11px',
                          fontWeight: 700,
                          color: isActive ? '#F59E0B' : 'oklch(0.50 0.015 240)',
                          textTransform: 'uppercase',
                          letterSpacing: '0.08em',
                        }}
                      >
                        {isActive ? 'Active Session' : 'Completed'}
                      </span>
                    </div>
                    <span
                      style={{
                        fontFamily: 'monospace',
                        fontSize: '12px',
                        color: 'oklch(0.55 0.015 240)',
                        background: 'oklch(0.22 0.025 245)',
                        padding: '2px 8px',
                        borderRadius: '6px',
                      }}
                    >
                      {getDuration(session.start_time, session.end_time)}
                    </span>
                  </div>

                  {/* Card body */}
                  <div className="p-5 space-y-4">
                    {/* Event name */}
                    <div>
                      <p
                        style={{
                          fontFamily: 'Space Grotesk, sans-serif',
                          fontSize: '16px',
                          fontWeight: 700,
                          color: 'oklch(0.97 0.005 240)',
                          marginBottom: '8px',
                        }}
                      >
                        {session.event_name || 'Event Session'}
                      </p>

                      {/* Date/time */}
                      <div className="flex items-center gap-2" style={{ marginBottom: '4px' }}>
                        <CalendarDays size={12} style={{ color: 'oklch(0.45 0.015 240)', flexShrink: 0 }} />
                        <span style={{ fontSize: '12px', color: 'oklch(0.55 0.015 240)' }}>
                          {formatDateTime(session.start_time)}
                          {session.end_time && ` → ${formatDateTime(session.end_time)}`}
                        </span>
                      </div>

                      {session.location && (
                        <div className="flex items-center gap-2">
                          <MapPin size={12} style={{ color: 'oklch(0.45 0.015 240)', flexShrink: 0 }} />
                          <span style={{ fontSize: '12px', color: 'oklch(0.55 0.015 240)' }}>
                            {session.location}
                          </span>
                        </div>
                      )}
                    </div>

                    {/* Stats row */}
                    <div
                      className="grid grid-cols-3 gap-3 pt-4"
                      style={{ borderTop: '1px solid oklch(0.22 0.025 245)' }}
                    >
                      <div
                        className="rounded-2xl p-3 text-center"
                        style={{ background: 'oklch(0.20 0.025 245)' }}
                      >
                        <Activity size={14} style={{ color: 'oklch(0.55 0.015 240)', margin: '0 auto 4px' }} />
                        <p
                          style={{
                            fontFamily: 'Space Grotesk, sans-serif',
                            fontSize: '20px',
                            fontWeight: 800,
                            color: 'oklch(0.97 0.005 240)',
                            lineHeight: 1,
                          }}
                        >
                          {total}
                        </p>
                        <p style={{ fontSize: '10px', color: 'oklch(0.50 0.015 240)', marginTop: '2px' }}>Total</p>
                      </div>
                      <div
                        className="rounded-2xl p-3 text-center"
                        style={{ background: 'rgba(16,185,129,0.08)', border: '1px solid rgba(16,185,129,0.15)' }}
                      >
                        <CheckCircle2 size={14} style={{ color: '#10B981', margin: '0 auto 4px' }} />
                        <p
                          style={{
                            fontFamily: 'Space Grotesk, sans-serif',
                            fontSize: '20px',
                            fontWeight: 800,
                            color: '#10B981',
                            lineHeight: 1,
                          }}
                        >
                          {valid}
                        </p>
                        <p style={{ fontSize: '10px', color: 'rgba(16,185,129,0.7)', marginTop: '2px' }}>Valid</p>
                      </div>
                      <div
                        className="rounded-2xl p-3 text-center"
                        style={{ background: 'rgba(244,63,94,0.08)', border: '1px solid rgba(244,63,94,0.15)' }}
                      >
                        <XCircle size={14} style={{ color: '#F43F5E', margin: '0 auto 4px' }} />
                        <p
                          style={{
                            fontFamily: 'Space Grotesk, sans-serif',
                            fontSize: '20px',
                            fontWeight: 800,
                            color: '#F43F5E',
                            lineHeight: 1,
                          }}
                        >
                          {invalid}
                        </p>
                        <p style={{ fontSize: '10px', color: 'rgba(244,63,94,0.7)', marginTop: '2px' }}>Invalid</p>
                      </div>
                    </div>

                    {/* Success rate bar */}
                    {total > 0 && (
                      <div>
                        <div className="flex items-center justify-between mb-1.5">
                          <span style={{ fontSize: '11px', color: 'oklch(0.50 0.015 240)' }}>Success Rate</span>
                          <span style={{ fontSize: '12px', fontWeight: 700, color: rate >= 80 ? '#10B981' : rate >= 50 ? '#F59E0B' : '#F43F5E' }}>
                            {rate}%
                          </span>
                        </div>
                        <div
                          className="w-full rounded-full overflow-hidden"
                          style={{ height: '4px', background: 'oklch(0.22 0.025 245)' }}
                        >
                          <div
                            className="h-full rounded-full"
                            style={{
                              width: `${rate}%`,
                              background: rate >= 80 ? '#10B981' : rate >= 50 ? '#F59E0B' : '#F43F5E',
                              transition: 'width 0.5s ease',
                            }}
                          />
                        </div>
                      </div>
                    )}

                    {/* CTA for active session */}
                    {isActive && (
                      <button
                        onClick={() => setLocation('/scan')}
                        className="w-full flex items-center justify-center gap-2 font-semibold"
                        style={{
                          height: '44px',
                          borderRadius: '14px',
                          background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)',
                          color: 'oklch(0.13 0.025 245)',
                          fontSize: '14px',
                          fontFamily: 'Space Grotesk, sans-serif',
                          border: 'none',
                          cursor: 'pointer',
                          marginTop: '4px',
                        }}
                      >
                        <Scan size={16} />
                        Continue Scanning
                      </button>
                    )}
                  </div>
                </div>
              );
            })}

            <p
              className="text-center pt-2"
              style={{ fontSize: '11px', color: 'oklch(0.40 0.015 240)' }}
            >
              Showing current session · Full history available in the admin portal
            </p>
          </div>
        )}
      </main>

      <style>{`
        @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
        @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
      `}</style>
    </div>
  );
}

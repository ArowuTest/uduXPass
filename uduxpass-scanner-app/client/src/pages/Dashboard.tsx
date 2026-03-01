/**
 * Dashboard.tsx — uduXPass Scanner PWA
 * Design: Enterprise Amber/Slate — "Professional Event Operations"
 * - Deep slate background with amber accent
 * - Real-time session timer, live stats
 * - Card-based control center layout
 * - Full API integration
 */

import { useState, useEffect, useCallback } from 'react';
import { useLocation } from 'wouter';
import { useAuth } from '@/contexts/AuthContext';
import { scannerApi } from '@/lib/api';
import { toast } from 'sonner';
import {
  Scan, LogOut, Plus, History, Activity, CheckCircle2,
  XCircle, TrendingUp, Calendar, Clock, StopCircle,
  Loader2, Wifi, WifiOff, ChevronRight, User, BarChart3,
  Zap, AlertCircle,
} from 'lucide-react';

interface Stats {
  total_sessions: number;
  total_scans: number;
  valid_scans: number;
  invalid_scans: number;
  success_rate: number;
  events_assigned: number;
  last_active_at: string;
}

interface ActiveSession {
  id: string;
  event_id: string;
  event_name?: string;
  location?: string;
  start_time: string;
  status: string;
  scans_count?: number;
}

function StatCard({
  icon, label, value, accent, sublabel,
}: {
  icon: React.ReactNode;
  label: string;
  value: string | number;
  accent?: string;
  sublabel?: string;
}) {
  return (
    <div
      className="rounded-2xl p-4 flex flex-col gap-2"
      style={{
        background: 'oklch(0.18 0.025 245)',
        border: '1px solid oklch(0.26 0.025 245)',
      }}
    >
      <div className="flex items-center gap-2">
        <div
          className="w-7 h-7 rounded-lg flex items-center justify-center"
          style={{ background: accent ? `${accent}18` : 'oklch(0.24 0.025 245)' }}
        >
          <span style={{ color: accent || 'oklch(0.60 0.015 240)' }}>{icon}</span>
        </div>
        <span style={{ fontSize: '11px', fontWeight: 600, color: 'oklch(0.55 0.015 240)', textTransform: 'uppercase', letterSpacing: '0.08em' }}>
          {label}
        </span>
      </div>
      <div>
        <p
          style={{
            fontFamily: 'Space Grotesk, sans-serif',
            fontSize: '1.75rem',
            fontWeight: 800,
            color: accent || 'oklch(0.97 0.005 240)',
            lineHeight: 1,
            letterSpacing: '-0.03em',
          }}
        >
          {value}
        </p>
        {sublabel && (
          <p style={{ fontSize: '11px', color: 'oklch(0.50 0.015 240)', marginTop: '2px' }}>{sublabel}</p>
        )}
      </div>
    </div>
  );
}

export default function Dashboard() {
  const [, setLocation] = useLocation();
  const { scanner, logout } = useAuth();
  const [stats, setStats] = useState<Stats | null>(null);
  const [activeSession, setActiveSession] = useState<ActiveSession | null>(null);
  const [isLoadingStats, setIsLoadingStats] = useState(true);
  const [isEndingSession, setIsEndingSession] = useState(false);
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [sessionDuration, setSessionDuration] = useState('');

  useEffect(() => {
    const on = () => setIsOnline(true);
    const off = () => setIsOnline(false);
    window.addEventListener('online', on);
    window.addEventListener('offline', off);
    return () => { window.removeEventListener('online', on); window.removeEventListener('offline', off); };
  }, []);

  const fetchData = useCallback(async () => {
    try {
      const [statsResp, sessionResp] = await Promise.allSettled([
        scannerApi.getStats(),
        scannerApi.getCurrentSession(),
      ]);
      if (statsResp.status === 'fulfilled') setStats(statsResp.value as Stats);
      if (sessionResp.status === 'fulfilled' && sessionResp.value) {
        setActiveSession(sessionResp.value as ActiveSession);
      } else {
        setActiveSession(null);
      }
    } catch (err) {
      console.error('Dashboard fetch error:', err);
    } finally {
      setIsLoadingStats(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, [fetchData]);

  // Live session timer
  useEffect(() => {
    if (!activeSession?.start_time) { setSessionDuration(''); return; }
    const update = () => {
      const diff = Math.floor((Date.now() - new Date(activeSession.start_time).getTime()) / 1000);
      const h = Math.floor(diff / 3600);
      const m = Math.floor((diff % 3600) / 60);
      const s = diff % 60;
      setSessionDuration(h > 0 ? `${h}h ${m}m ${s}s` : `${m}m ${s}s`);
    };
    update();
    const t = setInterval(update, 1000);
    return () => clearInterval(t);
  }, [activeSession?.start_time]);

  const handleEndSession = async () => {
    if (!activeSession) return;
    setIsEndingSession(true);
    try {
      await scannerApi.endSession();
      toast.success('Session ended');
      setActiveSession(null);
      fetchData();
    } catch (err: any) {
      toast.error(err?.response?.data?.message || 'Failed to end session');
    } finally {
      setIsEndingSession(false);
    }
  };

  const handleLogout = () => {
    logout();
    setLocation('/login');
  };

  const formatDate = (iso: string) => {
    try {
      return new Date(iso).toLocaleTimeString('en-NG', { hour: '2-digit', minute: '2-digit' });
    } catch { return iso; }
  };

  const successRate = stats ? Math.round(stats.success_rate) : 0;

  return (
    <div className="min-h-screen flex flex-col" style={{ background: 'oklch(0.13 0.025 245)' }}>

      {/* Header */}
      <header
        className="px-4 py-3 flex items-center justify-between sticky top-0 z-20"
        style={{
          background: 'oklch(0.16 0.025 245)',
          borderBottom: '1px solid oklch(0.24 0.025 245)',
          backdropFilter: 'blur(12px)',
        }}
      >
        <div className="flex items-center gap-3">
          {/* Avatar */}
          <div
            className="w-9 h-9 rounded-xl flex items-center justify-center"
            style={{ background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)' }}
          >
            <User size={16} color="oklch(0.13 0.025 245)" strokeWidth={2.5} />
          </div>
          <div>
            <p
              style={{
                fontFamily: 'Space Grotesk, sans-serif',
                fontSize: '14px',
                fontWeight: 700,
                color: 'oklch(0.97 0.005 240)',
                lineHeight: 1.2,
              }}
            >
              {scanner?.name || scanner?.username || 'Scanner'}
            </p>
            <p style={{ fontSize: '11px', color: 'oklch(0.50 0.015 240)', lineHeight: 1.2 }}>
              {scanner?.email || 'scanner@uduxpass.com'}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {/* Online indicator */}
          <div
            className="flex items-center gap-1.5 px-2.5 py-1 rounded-full"
            style={{
              background: isOnline ? 'rgba(16,185,129,0.12)' : 'rgba(244,63,94,0.12)',
              border: `1px solid ${isOnline ? 'rgba(16,185,129,0.3)' : 'rgba(244,63,94,0.3)'}`,
              color: isOnline ? '#10B981' : '#F43F5E',
              fontSize: '11px',
              fontWeight: 600,
            }}
          >
            {isOnline ? <Wifi size={11} /> : <WifiOff size={11} />}
            <span>{isOnline ? 'Online' : 'Offline'}</span>
          </div>

          {/* Logout */}
          <button
            onClick={handleLogout}
            className="w-8 h-8 rounded-lg flex items-center justify-center"
            style={{
              background: 'oklch(0.22 0.025 245)',
              border: '1px solid oklch(0.28 0.025 245)',
              color: 'oklch(0.55 0.015 240)',
              transition: 'all 150ms ease',
            }}
          >
            <LogOut size={15} />
          </button>
        </div>
      </header>

      {/* Main */}
      <main className="flex-1 px-4 py-5 space-y-4 max-w-lg mx-auto w-full pb-8">

        {/* Welcome line */}
        <div className="pt-1">
          <h1
            style={{
              fontFamily: 'Space Grotesk, sans-serif',
              fontSize: '1.5rem',
              fontWeight: 800,
              color: 'oklch(0.97 0.005 240)',
              letterSpacing: '-0.03em',
              lineHeight: 1.1,
            }}
          >
            {activeSession ? 'Session Active' : 'Ready to Scan'}
          </h1>
          <p style={{ fontSize: '13px', color: 'oklch(0.55 0.015 240)', marginTop: '4px' }}>
            {activeSession
              ? `Scanning for ${activeSession.event_name || 'event'}`
              : 'Start a session to begin validating tickets'}
          </p>
        </div>

        {/* Active Session Card */}
        {activeSession ? (
          <div
            className="rounded-2xl overflow-hidden"
            style={{
              background: 'linear-gradient(135deg, oklch(0.20 0.035 60) 0%, oklch(0.18 0.025 245) 100%)',
              border: '1px solid rgba(245,158,11,0.35)',
              boxShadow: '0 4px 24px rgba(245,158,11,0.12)',
            }}
          >
            {/* Top stripe */}
            <div
              className="h-1 w-full"
              style={{ background: 'linear-gradient(90deg, #F59E0B, #D97706)' }}
            />
            <div className="p-4">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-2">
                  <div
                    className="w-2 h-2 rounded-full"
                    style={{ background: '#F59E0B', boxShadow: '0 0 6px rgba(245,158,11,0.8)', animation: 'activePulse 1.5s ease-in-out infinite' }}
                  />
                  <span style={{ fontSize: '11px', fontWeight: 700, color: '#F59E0B', textTransform: 'uppercase', letterSpacing: '0.12em' }}>
                    Active Session
                  </span>
                </div>
                {sessionDuration && (
                  <span
                    className="font-mono"
                    style={{
                      fontSize: '12px',
                      fontWeight: 700,
                      color: '#F59E0B',
                      background: 'rgba(245,158,11,0.12)',
                      border: '1px solid rgba(245,158,11,0.25)',
                      padding: '2px 8px',
                      borderRadius: '20px',
                    }}
                  >
                    {sessionDuration}
                  </span>
                )}
              </div>

              <p
                style={{
                  fontFamily: 'Space Grotesk, sans-serif',
                  fontSize: '1.125rem',
                  fontWeight: 700,
                  color: 'oklch(0.97 0.005 240)',
                  marginBottom: '6px',
                }}
              >
                {activeSession.event_name || 'Event Session'}
              </p>

              <div className="flex items-center gap-1.5 mb-4" style={{ color: 'oklch(0.60 0.015 240)', fontSize: '12px' }}>
                <Clock size={12} />
                <span>Started at {formatDate(activeSession.start_time)}</span>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => setLocation('/scan')}
                  className="flex-1 flex items-center justify-center gap-2 h-11 rounded-xl font-semibold"
                  style={{
                    background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)',
                    color: 'oklch(0.13 0.025 245)',
                    fontSize: '14px',
                    fontFamily: 'Space Grotesk, sans-serif',
                    boxShadow: '0 4px 16px rgba(245,158,11,0.35)',
                    transition: 'all 150ms ease',
                  }}
                >
                  <Scan size={18} />
                  Scan Tickets
                </button>
                <button
                  onClick={handleEndSession}
                  disabled={isEndingSession}
                  className="w-11 h-11 rounded-xl flex items-center justify-center"
                  style={{
                    background: 'rgba(244,63,94,0.12)',
                    border: '1px solid rgba(244,63,94,0.3)',
                    color: '#F43F5E',
                    transition: 'all 150ms ease',
                  }}
                >
                  {isEndingSession ? <Loader2 size={18} className="animate-spin" /> : <StopCircle size={18} />}
                </button>
              </div>
            </div>
          </div>
        ) : (
          /* No session CTA */
          <div
            className="rounded-2xl p-5 text-center"
            style={{
              background: 'oklch(0.18 0.025 245)',
              border: '1px solid oklch(0.26 0.025 245)',
              borderStyle: 'dashed',
            }}
          >
            <div
              className="w-12 h-12 rounded-2xl flex items-center justify-center mx-auto mb-3"
              style={{ background: 'oklch(0.22 0.025 245)' }}
            >
              <Scan size={22} style={{ color: 'oklch(0.50 0.015 240)' }} />
            </div>
            <p
              style={{
                fontFamily: 'Space Grotesk, sans-serif',
                fontSize: '15px',
                fontWeight: 700,
                color: 'oklch(0.97 0.005 240)',
                marginBottom: '4px',
              }}
            >
              No Active Session
            </p>
            <p style={{ fontSize: '13px', color: 'oklch(0.55 0.015 240)', marginBottom: '16px' }}>
              Start a session to begin validating tickets
            </p>
            <button
              onClick={() => setLocation('/create-session')}
              className="flex items-center justify-center gap-2 h-11 px-6 rounded-xl font-semibold mx-auto"
              style={{
                background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)',
                color: 'oklch(0.13 0.025 245)',
                fontSize: '14px',
                fontFamily: 'Space Grotesk, sans-serif',
                boxShadow: '0 4px 16px rgba(245,158,11,0.35)',
              }}
            >
              <Plus size={18} />
              Start New Session
            </button>
          </div>
        )}

        {/* Stats Section */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <BarChart3 size={14} style={{ color: 'oklch(0.50 0.015 240)' }} />
            <span
              style={{
                fontSize: '11px',
                fontWeight: 700,
                color: 'oklch(0.50 0.015 240)',
                textTransform: 'uppercase',
                letterSpacing: '0.12em',
              }}
            >
              Your Statistics
            </span>
          </div>

          {isLoadingStats ? (
            <div className="grid grid-cols-2 gap-3">
              {[...Array(4)].map((_, i) => (
                <div
                  key={i}
                  className="rounded-2xl p-4 animate-pulse"
                  style={{ background: 'oklch(0.18 0.025 245)', height: '96px' }}
                />
              ))}
            </div>
          ) : (
            <>
              <div className="grid grid-cols-2 gap-3 mb-3">
                <StatCard
                  icon={<Activity size={14} />}
                  label="Total Scans"
                  value={stats?.total_scans ?? 0}
                  sublabel="all time"
                />
                <StatCard
                  icon={<CheckCircle2 size={14} />}
                  label="Valid"
                  value={stats?.valid_scans ?? 0}
                  accent="#10B981"
                  sublabel="tickets passed"
                />
                <StatCard
                  icon={<XCircle size={14} />}
                  label="Invalid"
                  value={stats?.invalid_scans ?? 0}
                  accent="#F43F5E"
                  sublabel="rejected"
                />
                <StatCard
                  icon={<TrendingUp size={14} />}
                  label="Success Rate"
                  value={`${successRate}%`}
                  accent={successRate >= 80 ? '#10B981' : successRate >= 50 ? '#F59E0B' : '#F43F5E'}
                />
              </div>

              {/* Secondary stats */}
              <div className="grid grid-cols-2 gap-3">
                <div
                  className="rounded-2xl p-3 flex items-center gap-3"
                  style={{ background: 'oklch(0.18 0.025 245)', border: '1px solid oklch(0.26 0.025 245)' }}
                >
                  <div
                    className="w-8 h-8 rounded-xl flex items-center justify-center flex-shrink-0"
                    style={{ background: 'rgba(245,158,11,0.12)' }}
                  >
                    <Calendar size={15} style={{ color: '#F59E0B' }} />
                  </div>
                  <div>
                    <p style={{ fontSize: '11px', color: 'oklch(0.50 0.015 240)' }}>Events Assigned</p>
                    <p style={{ fontFamily: 'Space Grotesk, sans-serif', fontWeight: 700, color: 'oklch(0.97 0.005 240)', fontSize: '18px', lineHeight: 1.2 }}>
                      {stats?.events_assigned ?? 0}
                    </p>
                  </div>
                </div>
                <div
                  className="rounded-2xl p-3 flex items-center gap-3"
                  style={{ background: 'oklch(0.18 0.025 245)', border: '1px solid oklch(0.26 0.025 245)' }}
                >
                  <div
                    className="w-8 h-8 rounded-xl flex items-center justify-center flex-shrink-0"
                    style={{ background: 'oklch(0.24 0.025 245)' }}
                  >
                    <History size={15} style={{ color: 'oklch(0.55 0.015 240)' }} />
                  </div>
                  <div>
                    <p style={{ fontSize: '11px', color: 'oklch(0.50 0.015 240)' }}>Total Sessions</p>
                    <p style={{ fontFamily: 'Space Grotesk, sans-serif', fontWeight: 700, color: 'oklch(0.97 0.005 240)', fontSize: '18px', lineHeight: 1.2 }}>
                      {stats?.total_sessions ?? 0}
                    </p>
                  </div>
                </div>
              </div>
            </>
          )}
        </div>

        {/* Quick Actions */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Zap size={14} style={{ color: 'oklch(0.50 0.015 240)' }} />
            <span
              style={{
                fontSize: '11px',
                fontWeight: 700,
                color: 'oklch(0.50 0.015 240)',
                textTransform: 'uppercase',
                letterSpacing: '0.12em',
              }}
            >
              Quick Actions
            </span>
          </div>

          <div className="space-y-2">
            {activeSession ? (
              <button
                onClick={() => setLocation('/scan')}
                className="w-full flex items-center justify-between p-4 rounded-2xl group"
                style={{
                  background: 'oklch(0.18 0.025 245)',
                  border: '1px solid oklch(0.26 0.025 245)',
                  transition: 'all 150ms ease',
                }}
              >
                <div className="flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-xl flex items-center justify-center"
                    style={{ background: 'rgba(245,158,11,0.12)' }}
                  >
                    <Scan size={18} style={{ color: '#F59E0B' }} />
                  </div>
                  <div className="text-left">
                    <p style={{ fontSize: '14px', fontWeight: 600, color: 'oklch(0.97 0.005 240)' }}>Scan Tickets</p>
                    <p style={{ fontSize: '12px', color: 'oklch(0.55 0.015 240)' }}>Camera or manual entry</p>
                  </div>
                </div>
                <ChevronRight size={16} style={{ color: 'oklch(0.45 0.015 240)' }} />
              </button>
            ) : (
              <button
                onClick={() => setLocation('/create-session')}
                className="w-full flex items-center justify-between p-4 rounded-2xl"
                style={{
                  background: 'oklch(0.18 0.025 245)',
                  border: '1px solid oklch(0.26 0.025 245)',
                  transition: 'all 150ms ease',
                }}
              >
                <div className="flex items-center gap-3">
                  <div
                    className="w-10 h-10 rounded-xl flex items-center justify-center"
                    style={{ background: 'rgba(245,158,11,0.12)' }}
                  >
                    <Plus size={18} style={{ color: '#F59E0B' }} />
                  </div>
                  <div className="text-left">
                    <p style={{ fontSize: '14px', fontWeight: 600, color: 'oklch(0.97 0.005 240)' }}>New Session</p>
                    <p style={{ fontSize: '12px', color: 'oklch(0.55 0.015 240)' }}>Start scanning for an event</p>
                  </div>
                </div>
                <ChevronRight size={16} style={{ color: 'oklch(0.45 0.015 240)' }} />
              </button>
            )}

            <button
              onClick={() => setLocation('/history')}
              className="w-full flex items-center justify-between p-4 rounded-2xl"
              style={{
                background: 'oklch(0.18 0.025 245)',
                border: '1px solid oklch(0.26 0.025 245)',
                transition: 'all 150ms ease',
              }}
            >
              <div className="flex items-center gap-3">
                <div
                  className="w-10 h-10 rounded-xl flex items-center justify-center"
                  style={{ background: 'oklch(0.24 0.025 245)' }}
                >
                  <History size={18} style={{ color: 'oklch(0.55 0.015 240)' }} />
                </div>
                <div className="text-left">
                  <p style={{ fontSize: '14px', fontWeight: 600, color: 'oklch(0.97 0.005 240)' }}>Session History</p>
                  <p style={{ fontSize: '12px', color: 'oklch(0.55 0.015 240)' }}>View past scanning sessions</p>
                </div>
              </div>
              <ChevronRight size={16} style={{ color: 'oklch(0.45 0.015 240)' }} />
            </button>
          </div>
        </div>

        {/* Offline warning */}
        {!isOnline && (
          <div
            className="flex items-center gap-3 p-3 rounded-xl"
            style={{
              background: 'rgba(245,158,11,0.08)',
              border: '1px solid rgba(245,158,11,0.25)',
            }}
          >
            <AlertCircle size={16} style={{ color: '#F59E0B', flexShrink: 0 }} />
            <p style={{ fontSize: '12px', color: '#F59E0B' }}>
              You are offline. Scans will be queued and synced when connection is restored.
            </p>
          </div>
        )}
      </main>

      <style>{`
        @keyframes activePulse {
          0%, 100% { opacity: 1; box-shadow: 0 0 6px rgba(245,158,11,0.8); }
          50% { opacity: 0.5; box-shadow: 0 0 12px rgba(245,158,11,0.4); }
        }
      `}</style>
    </div>
  );
}

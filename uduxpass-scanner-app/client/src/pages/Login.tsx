/**
 * Login.tsx — uduXPass Scanner PWA
 * Design: Enterprise Amber/Slate — "Professional Event Operations"
 * - Deep slate navy background with ambient amber glow
 * - Space Grotesk display font for brand authority
 * - Large touch targets (56px CTA), smooth transitions
 * - Real API integration via AuthContext
 */

import { useState, useRef, useEffect } from 'react';
import { useLocation } from 'wouter';
import { useAuth } from '@/contexts/AuthContext';
import { Eye, EyeOff, Scan, Shield, Zap, QrCode } from 'lucide-react';
import { toast } from 'sonner';

export default function Login() {
  const [, setLocation] = useLocation();
  const { login, scanner } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const usernameRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (scanner) setLocation('/dashboard');
  }, [scanner, setLocation]);

  useEffect(() => {
    usernameRef.current?.focus();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      setError('Please enter your username and password.');
      return;
    }
    setIsLoading(true);
    setError('');
    try {
      await login({ username: username.trim(), password });
      toast.success('Welcome back!');
      setLocation('/dashboard');
    } catch (err: any) {
      const msg = err?.response?.data?.message || err?.message || 'Invalid credentials. Please try again.';
      setError(msg);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex flex-col"
      style={{ background: 'oklch(0.13 0.025 245)' }}
    >
      {/* Ambient background glows */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none select-none">
        <div
          className="absolute -top-32 -right-32 w-80 h-80 rounded-full"
          style={{ background: 'radial-gradient(circle, rgba(245,158,11,0.12) 0%, transparent 70%)' }}
        />
        <div
          className="absolute -bottom-32 -left-32 w-72 h-72 rounded-full"
          style={{ background: 'radial-gradient(circle, rgba(59,130,246,0.08) 0%, transparent 70%)' }}
        />
        {/* Subtle grid */}
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `linear-gradient(oklch(0.22 0.025 245 / 0.4) 1px, transparent 1px),
                              linear-gradient(90deg, oklch(0.22 0.025 245 / 0.4) 1px, transparent 1px)`,
            backgroundSize: '48px 48px',
            opacity: 0.5,
          }}
        />
      </div>

      {/* Content */}
      <div className="flex-1 flex flex-col items-center justify-center px-6 py-12 relative z-10">
        <div className="w-full max-w-sm">

          {/* Brand Header */}
          <div className="text-center mb-10">
            <div className="flex items-center justify-center mb-5">
              <div className="relative">
                <div
                  className="w-20 h-20 rounded-[22px] flex items-center justify-center"
                  style={{
                    background: 'linear-gradient(135deg, #F59E0B 0%, #D97706 100%)',
                    boxShadow: '0 8px 32px rgba(245,158,11,0.45), 0 2px 8px rgba(0,0,0,0.3)',
                  }}
                >
                  <Scan size={38} color="oklch(0.13 0.025 245)" strokeWidth={2.5} />
                </div>
                {/* Animated ring */}
                <div
                  className="absolute inset-0 rounded-[22px]"
                  style={{
                    border: '2px solid rgba(245,158,11,0.4)',
                    animation: 'loginPulse 2.5s ease-out infinite',
                  }}
                />
              </div>
            </div>

            <h1
              style={{
                fontFamily: 'Space Grotesk, sans-serif',
                fontSize: '2rem',
                fontWeight: 800,
                letterSpacing: '-0.04em',
                color: 'oklch(0.97 0.005 240)',
                lineHeight: 1,
                marginBottom: '8px',
              }}
            >
              uduX<span style={{ color: '#F59E0B' }}>Pass</span>
            </h1>
            <p
              style={{
                fontSize: '11px',
                fontWeight: 600,
                letterSpacing: '0.18em',
                textTransform: 'uppercase',
                color: 'oklch(0.50 0.015 240)',
              }}
            >
              Scanner Portal
            </p>
          </div>

          {/* Login Card */}
          <div
            className="rounded-2xl p-6 mb-5"
            style={{
              background: 'oklch(0.18 0.025 245)',
              border: '1px solid oklch(0.26 0.025 245)',
              boxShadow: '0 24px 64px rgba(0,0,0,0.5)',
            }}
          >
            <div className="mb-5">
              <h2
                style={{
                  fontFamily: 'Space Grotesk, sans-serif',
                  fontSize: '1.25rem',
                  fontWeight: 700,
                  color: 'oklch(0.97 0.005 240)',
                  marginBottom: '4px',
                }}
              >
                Sign In
              </h2>
              <p style={{ fontSize: '14px', color: 'oklch(0.55 0.015 240)' }}>
                Enter your scanner credentials to continue
              </p>
            </div>

            {/* Error message */}
            {error && (
              <div
                className="flex items-start gap-2.5 p-3 rounded-xl mb-4 text-sm fade-in"
                style={{
                  background: 'rgba(244,63,94,0.1)',
                  border: '1px solid rgba(244,63,94,0.25)',
                  color: '#FB7185',
                }}
              >
                <Shield size={15} className="mt-0.5 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              {/* Username */}
              <div>
                <label
                  className="block mb-2"
                  style={{
                    fontSize: '11px',
                    fontWeight: 700,
                    letterSpacing: '0.1em',
                    textTransform: 'uppercase',
                    color: 'oklch(0.55 0.015 240)',
                  }}
                >
                  Username
                </label>
                <input
                  ref={usernameRef}
                  type="text"
                  value={username}
                  onChange={(e) => { setUsername(e.target.value); setError(''); }}
                  placeholder="e.g. scanner1"
                  autoComplete="username"
                  autoCapitalize="none"
                  autoCorrect="off"
                  spellCheck={false}
                  className="input-field"
                  disabled={isLoading}
                />
              </div>

              {/* Password */}
              <div>
                <label
                  className="block mb-2"
                  style={{
                    fontSize: '11px',
                    fontWeight: 700,
                    letterSpacing: '0.1em',
                    textTransform: 'uppercase',
                    color: 'oklch(0.55 0.015 240)',
                  }}
                >
                  Password
                </label>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => { setPassword(e.target.value); setError(''); }}
                    placeholder="Enter your password"
                    autoComplete="current-password"
                    className="input-field"
                    style={{ paddingRight: '48px' }}
                    disabled={isLoading}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 p-1.5 rounded-lg"
                    style={{ color: 'oklch(0.50 0.015 240)', transition: 'color 150ms ease' }}
                    tabIndex={-1}
                  >
                    {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                  </button>
                </div>
              </div>

              {/* Submit Button */}
              <div className="pt-1">
                <button
                  type="submit"
                  className="btn-brand"
                  disabled={isLoading || !username.trim() || !password.trim()}
                >
                  {isLoading ? (
                    <>
                      <div className="spinner" style={{ width: 20, height: 20, borderTopColor: 'oklch(0.13 0.025 245)', borderColor: 'rgba(0,0,0,0.2)' }} />
                      <span>Signing in...</span>
                    </>
                  ) : (
                    <>
                      <QrCode size={20} />
                      <span>Sign In to Scanner</span>
                    </>
                  )}
                </button>
              </div>
            </form>
          </div>

          {/* Feature Pills */}
          <div className="flex items-center justify-center gap-2 flex-wrap">
            {[
              { icon: <Zap size={11} />, label: 'Real-time' },
              { icon: <Shield size={11} />, label: 'Secure' },
              { icon: <QrCode size={11} />, label: 'QR + Manual' },
            ].map((f) => (
              <div
                key={f.label}
                className="flex items-center gap-1.5 px-3 py-1.5 rounded-full"
                style={{
                  background: 'oklch(0.18 0.025 245)',
                  border: '1px solid oklch(0.26 0.025 245)',
                  color: 'oklch(0.55 0.015 240)',
                  fontSize: '11px',
                  fontWeight: 500,
                }}
              >
                <span style={{ color: '#F59E0B' }}>{f.icon}</span>
                {f.label}
              </div>
            ))}
          </div>

          <p
            className="text-center mt-8"
            style={{ fontSize: '11px', color: 'oklch(0.40 0.015 240)' }}
          >
            uduXPass &copy; {new Date().getFullYear()} — Professional Event Management
          </p>
        </div>
      </div>

      <style>{`
        @keyframes loginPulse {
          0% { transform: scale(1); opacity: 0.7; }
          70% { transform: scale(1.4); opacity: 0; }
          100% { transform: scale(1.4); opacity: 0; }
        }
      `}</style>
    </div>
  );
}

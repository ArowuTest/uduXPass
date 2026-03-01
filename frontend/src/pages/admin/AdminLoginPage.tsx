/*
 * AdminLoginPage — uduXPass Design System
 * Dark navy, amber accent, Syne headings — enterprise admin portal feel
 */
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Eye, EyeOff, Shield, Lock, Mail, AlertCircle } from 'lucide-react'

const AdminLoginPage: React.FC = () => {
  const [email, setEmail] = useState('admin@uduxpass.com')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const { adminLogin } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email || !password) { setError('Please fill in all fields'); return }
    setIsLoading(true)
    setError('')
    try {
      await adminLogin(email, password)
      navigate('/admin/dashboard')
    } catch {
      setError('Invalid credentials. Please check your email and password.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex" style={{ background: 'var(--brand-navy)' }}>
      {/* Left panel — branding */}
      <div className="hidden lg:flex flex-col justify-between w-2/5 p-12 relative overflow-hidden"
        style={{ background: 'var(--brand-surface)', borderRight: '1px solid rgba(255,255,255,0.07)' }}>
        <div className="absolute inset-0 opacity-5"
          style={{ backgroundImage: 'radial-gradient(circle at 20% 50%, var(--brand-amber) 0%, transparent 60%)' }} />
        <div className="relative z-10">
          <div className="flex items-center gap-3 mb-16">
            <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ background: 'var(--brand-amber)' }}>
              <Shield className="w-5 h-5" style={{ color: '#0f1729' }} />
            </div>
            <span className="text-xl font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
              uduX<span style={{ color: 'var(--brand-amber)' }}>Admin</span>
            </span>
          </div>
          <h2 className="text-4xl font-bold mb-4" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9', lineHeight: 1.15 }}>
            Manage your<br />events with<br /><span style={{ color: 'var(--brand-amber)' }}>confidence.</span>
          </h2>
          <p className="text-base" style={{ color: '#64748b', lineHeight: 1.7 }}>
            The uduXPass admin portal gives you full control over events, orders, tickets, and analytics.
          </p>
        </div>
        <div className="relative z-10">
          {[
            { label: 'Role-based access control', desc: 'Fine-grained permissions per admin' },
            { label: 'Real-time analytics', desc: 'Live revenue and ticket data' },
            { label: 'Scanner management', desc: 'Manage event entry scanners' },
          ].map(item => (
            <div key={item.label} className="flex items-start gap-3 mb-4">
              <div className="w-1.5 h-1.5 rounded-full mt-2 flex-shrink-0" style={{ background: 'var(--brand-amber)' }} />
              <div>
                <p className="text-sm font-semibold" style={{ color: '#f1f5f9', fontFamily: 'var(--font-display)' }}>{item.label}</p>
                <p className="text-xs" style={{ color: '#64748b' }}>{item.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Right panel — form */}
      <div className="flex-1 flex items-center justify-center px-6 py-12">
        <div className="w-full max-w-sm">
          {/* Mobile logo */}
          <div className="flex items-center gap-2 mb-10 lg:hidden">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: 'var(--brand-amber)' }}>
              <Shield className="w-4 h-4" style={{ color: '#0f1729' }} />
            </div>
            <span className="text-lg font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
              uduX<span style={{ color: 'var(--brand-amber)' }}>Admin</span>
            </span>
          </div>

          <div className="mb-8">
            <h1 className="text-2xl font-bold mb-1" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
              Admin Sign In
            </h1>
            <p className="text-sm" style={{ color: '#64748b' }}>Enter your credentials to access the admin portal</p>
          </div>

          {error && (
            <div className="flex items-start gap-3 p-4 rounded-xl mb-6"
              style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)' }}>
              <AlertCircle className="w-4 h-4 flex-shrink-0 mt-0.5" style={{ color: '#f87171' }} />
              <p className="text-sm" style={{ color: '#f87171' }}>{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <Label htmlFor="email" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>
                Email Address
              </Label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={e => { setEmail(e.target.value); setError('') }}
                  placeholder="admin@uduxpass.com"
                  className="pl-10"
                  style={{
                    background: 'var(--brand-surface)',
                    border: '1px solid rgba(255,255,255,0.1)',
                    color: '#f1f5f9',
                    height: '44px',
                  }}
                  required
                />
              </div>
            </div>

            <div>
              <Label htmlFor="password" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>
                Password
              </Label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={e => { setPassword(e.target.value); setError('') }}
                  placeholder="Enter your password"
                  className="pl-10 pr-10"
                  style={{
                    background: 'var(--brand-surface)',
                    border: '1px solid rgba(255,255,255,0.1)',
                    color: '#f1f5f9',
                    height: '44px',
                  }}
                  required
                />
                <button type="button" onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2" style={{ color: '#475569' }}>
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
            </div>

            <Button type="submit" disabled={isLoading} className="w-full h-11 font-bold text-base"
              style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)', marginTop: '8px' }}>
              {isLoading ? (
                <div className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-t-transparent rounded-full animate-spin" style={{ borderColor: '#0f1729', borderTopColor: 'transparent' }} />
                  Signing In...
                </div>
              ) : 'Sign In to Admin Portal'}
            </Button>
          </form>

          <p className="text-center text-xs mt-8" style={{ color: '#475569' }}>
            This portal is restricted to authorized administrators only.
          </p>
        </div>
      </div>
    </div>
  )
}

export default AdminLoginPage

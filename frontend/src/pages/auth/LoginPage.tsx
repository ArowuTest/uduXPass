/*
 * LoginPage — uduXPass Design System
 * Dark navy, amber accent, Syne headings
 */
import React, { useState } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '../../contexts/AuthContext'
import { Eye, EyeOff, Mail, Lock, Ticket, AlertCircle, ArrowLeft } from 'lucide-react'
import { toast } from '@/components/ui/toaster'

interface LocationState { from?: { pathname: string } }

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const from = (location.state as LocationState)?.from?.pathname || '/home'

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email || !password) { setError('Please fill in all fields'); return }
    setIsLoading(true)
    setError('')
    try {
      await login(email, password)
      toast.success('Welcome back!')
      navigate(from, { replace: true })
    } catch (err: any) {
      setError(err?.message || 'Invalid email or password. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-12" style={{ background: 'var(--brand-navy)' }}>
      <div className="w-full max-w-md">
        {/* Back */}
        <Link to="/home" className="inline-flex items-center gap-2 text-sm mb-8"
          style={{ color: '#64748b', textDecoration: 'none' }}>
          <ArrowLeft className="w-4 h-4" /> Back to Home
        </Link>

        {/* Logo */}
        <div className="flex items-center gap-2 mb-8">
          <div className="w-9 h-9 rounded-xl flex items-center justify-center" style={{ background: 'var(--brand-amber)' }}>
            <Ticket className="w-5 h-5" style={{ color: '#0f1729' }} />
          </div>
          <span className="text-xl font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
            uduX<span style={{ color: 'var(--brand-amber)' }}>Pass</span>
          </span>
        </div>

        <div className="p-8 rounded-2xl" style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.08)' }}>
          <h1 className="text-2xl font-bold mb-1" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Sign In</h1>
          <p className="text-sm mb-6" style={{ color: '#64748b' }}>
            Don't have an account? <Link to="/register" style={{ color: 'var(--brand-amber)', textDecoration: 'none', fontWeight: 600 }}>Create one</Link>
          </p>

          {error && (
            <div className="flex items-start gap-3 p-4 rounded-xl mb-5"
              style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)' }}>
              <AlertCircle className="w-4 h-4 flex-shrink-0 mt-0.5" style={{ color: '#f87171' }} />
              <p className="text-sm" style={{ color: '#f87171' }}>{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <Label htmlFor="email" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>Email</Label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input id="email" type="email" value={email} onChange={e => { setEmail(e.target.value); setError('') }}
                  placeholder="you@example.com" className="pl-10 h-11"
                  style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} required />
              </div>
            </div>
            <div>
              <div className="flex items-center justify-between mb-1.5">
                <Label htmlFor="password" className="text-sm font-medium" style={{ color: '#94a3b8' }}>Password</Label>
                <Link to="/forgot-password" className="text-xs" style={{ color: 'var(--brand-amber)', textDecoration: 'none' }}>Forgot password?</Link>
              </div>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input id="password" type={showPassword ? 'text' : 'password'} value={password}
                  onChange={e => { setPassword(e.target.value); setError('') }}
                  placeholder="Enter your password" className="pl-10 pr-10 h-11"
                  style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} required />
                <button type="button" onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2" style={{ color: '#475569' }}>
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
            </div>
            <Button type="submit" disabled={isLoading} className="w-full h-11 font-bold text-base mt-2"
              style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)' }}>
              {isLoading ? (
                <div className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-t-transparent rounded-full animate-spin" style={{ borderColor: '#0f1729', borderTopColor: 'transparent' }} />
                  Signing In...
                </div>
              ) : 'Sign In'}
            </Button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default LoginPage

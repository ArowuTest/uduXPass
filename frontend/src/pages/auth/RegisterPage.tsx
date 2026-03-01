/*
 * RegisterPage — uduXPass Design System
 * Dark navy, amber accent, Syne headings
 */
import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuth } from '../../contexts/AuthContext'
import { Eye, EyeOff, Mail, Lock, User, Phone, Ticket, AlertCircle, ArrowLeft, CheckCircle } from 'lucide-react'
import { toast } from '@/components/ui/toaster'

const RegisterPage: React.FC = () => {
  const [form, setForm] = useState({ firstName: '', lastName: '', email: '', phone: '', password: '' })
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const { register } = useAuth()
  const navigate = useNavigate()

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm(prev => ({ ...prev, [field]: e.target.value }))
    setError('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const { firstName, lastName, email, phone, password } = form
    if (!firstName || !lastName || !email || !password) { setError('Please fill in all required fields'); return }
    if (password.length < 8) { setError('Password must be at least 8 characters'); return }
    setIsLoading(true)
    setError('')
    try {
      await register({ firstName, lastName, email, phone, password })
      toast.success('Account created! Welcome to uduXPass.')
      navigate('/home')
    } catch (err: any) {
      setError(err?.message || 'Registration failed. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const passwordStrength = form.password.length === 0 ? 0 : form.password.length < 6 ? 1 : form.password.length < 10 ? 2 : 3
  const strengthColors = ['', '#ef4444', '#F59E0B', '#10b981']
  const strengthLabels = ['', 'Weak', 'Fair', 'Strong']

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
          <h1 className="text-2xl font-bold mb-1" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Create Account</h1>
          <p className="text-sm mb-6" style={{ color: '#64748b' }}>
            Already have an account? <Link to="/login" style={{ color: 'var(--brand-amber)', textDecoration: 'none', fontWeight: 600 }}>Sign in</Link>
          </p>

          {error && (
            <div className="flex items-start gap-3 p-4 rounded-xl mb-5"
              style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)' }}>
              <AlertCircle className="w-4 h-4 flex-shrink-0 mt-0.5" style={{ color: '#f87171' }} />
              <p className="text-sm" style={{ color: '#f87171' }}>{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-3">
              <div>
                <Label htmlFor="firstName" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>First Name *</Label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                  <Input id="firstName" value={form.firstName} onChange={set('firstName')} placeholder="John" className="pl-10 h-11"
                    style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} required />
                </div>
              </div>
              <div>
                <Label htmlFor="lastName" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>Last Name *</Label>
                <div className="relative">
                  <User className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                  <Input id="lastName" value={form.lastName} onChange={set('lastName')} placeholder="Doe" className="pl-10 h-11"
                    style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} required />
                </div>
              </div>
            </div>

            <div>
              <Label htmlFor="email" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>Email *</Label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input id="email" type="email" value={form.email} onChange={set('email')} placeholder="you@example.com" className="pl-10 h-11"
                  style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} required />
              </div>
            </div>

            <div>
              <Label htmlFor="phone" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>Phone Number</Label>
              <div className="relative">
                <Phone className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input id="phone" type="tel" value={form.phone} onChange={set('phone')} placeholder="+234 800 000 0000" className="pl-10 h-11"
                  style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} />
              </div>
            </div>

            <div>
              <Label htmlFor="password" className="text-sm font-medium mb-1.5 block" style={{ color: '#94a3b8' }}>Password *</Label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
                <Input id="password" type={showPassword ? 'text' : 'password'} value={form.password} onChange={set('password')}
                  placeholder="Min. 8 characters" className="pl-10 pr-10 h-11"
                  style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }} required />
                <button type="button" onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2" style={{ color: '#475569' }}>
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
              {form.password && (
                <div className="flex items-center gap-2 mt-2">
                  <div className="flex gap-1 flex-1">
                    {[1,2,3].map(i => (
                      <div key={i} className="h-1 flex-1 rounded-full transition-all duration-300"
                        style={{ background: i <= passwordStrength ? strengthColors[passwordStrength] : 'rgba(255,255,255,0.1)' }} />
                    ))}
                  </div>
                  <span className="text-xs" style={{ color: strengthColors[passwordStrength] }}>{strengthLabels[passwordStrength]}</span>
                </div>
              )}
            </div>

            <Button type="submit" disabled={isLoading} className="w-full h-11 font-bold text-base mt-2"
              style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)' }}>
              {isLoading ? (
                <div className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-t-transparent rounded-full animate-spin" style={{ borderColor: '#0f1729', borderTopColor: 'transparent' }} />
                  Creating Account...
                </div>
              ) : 'Create Account'}
            </Button>
          </form>

          <p className="text-xs text-center mt-5" style={{ color: '#475569' }}>
            By creating an account, you agree to our{' '}
            <Link to="/terms" style={{ color: 'var(--brand-amber)', textDecoration: 'none' }}>Terms of Service</Link>
            {' '}and{' '}
            <Link to="/privacy" style={{ color: 'var(--brand-amber)', textDecoration: 'none' }}>Privacy Policy</Link>
          </p>
        </div>
      </div>
    </div>
  )
}

export default RegisterPage

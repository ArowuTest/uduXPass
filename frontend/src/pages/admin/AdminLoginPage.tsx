/*
 * AdminLoginPage — uduXPass Admin Portal
 * Design: Dark navy (#0f1729) + amber (#F59E0B) brand system
 * Split-panel layout: feature highlights left, login form right
 */
import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Eye, EyeOff, Shield, ArrowLeft, BarChart3, Users, Ticket, Zap } from 'lucide-react'

interface LoginFormData {
  email: string
  password: string
}

const AdminLoginPage: React.FC = () => {
  const [formData, setFormData] = useState<LoginFormData>({
    email: 'admin@uduxpass.com',
    password: ''
  })
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const { adminLogin } = useAuth()
  const navigate = useNavigate()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
    if (error) setError('')
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!formData.email || !formData.password) {
      setError('Please fill in all fields')
      return
    }
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
    setIsLoading(false)
  }

  const features = [
    { icon: BarChart3, title: 'Real-time Analytics', desc: 'Live revenue, ticket sales, and attendance dashboards' },
    { icon: Users,    title: 'User Management',    desc: 'Manage customers, organizers, and admin roles' },
    { icon: Ticket,   title: 'Event Control',      desc: 'Create events, manage tiers, and track inventory' },
    { icon: Zap,      title: 'Scanner Integration', desc: 'Monitor live gate scanning and entry validation' },
  ]

  return (
    <div className="min-h-screen flex" style={{ background: '#0f1729' }}>
      {/* Left panel — feature highlights */}
      <div className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12"
        style={{ background: 'linear-gradient(135deg, #0f1729 0%, #1a2744 100%)', borderRight: '1px solid rgba(245,158,11,0.15)' }}>
        <div>
          <Link to="/" className="inline-flex items-center gap-2 text-amber-400 hover:text-amber-300 transition-colors text-sm font-medium">
            <ArrowLeft className="w-4 h-4" />
            Back to uduXPass
          </Link>
        </div>

        <div>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ background: '#F59E0B' }}>
              <Shield className="w-5 h-5 text-slate-900" />
            </div>
            <div>
              <div className="text-white font-bold text-lg leading-none">uduXPass</div>
              <div className="text-amber-400 text-xs font-medium tracking-widest uppercase">Admin Portal</div>
            </div>
          </div>

          <h1 className="text-4xl font-bold text-white mb-4 leading-tight">
            Manage your events<br />
            <span style={{ color: '#F59E0B' }}>with full control.</span>
          </h1>
          <p className="text-slate-400 text-base leading-relaxed mb-10">
            The uduXPass administration platform gives you complete oversight of every event, ticket, order, and attendee on the platform.
          </p>

          <div className="space-y-5">
            {features.map(({ icon: Icon, title, desc }) => (
              <div key={title} className="flex items-start gap-4">
                <div className="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0"
                  style={{ background: 'rgba(245,158,11,0.12)', border: '1px solid rgba(245,158,11,0.25)' }}>
                  <Icon className="w-4 h-4" style={{ color: '#F59E0B' }} />
                </div>
                <div>
                  <div className="text-white text-sm font-semibold">{title}</div>
                  <div className="text-slate-400 text-xs mt-0.5">{desc}</div>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="text-slate-600 text-xs">
          © 2026 uduXPass. All rights reserved.
        </div>
      </div>

      {/* Right panel — login form */}
      <div className="w-full lg:w-1/2 flex flex-col items-center justify-center p-8">
        {/* Mobile back link */}
        <div className="lg:hidden w-full max-w-md mb-6">
          <Link to="/" className="inline-flex items-center gap-2 text-amber-400 hover:text-amber-300 transition-colors text-sm">
            <ArrowLeft className="w-4 h-4" />
            Back to uduXPass
          </Link>
        </div>

        <div className="w-full max-w-md">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-white mb-1">Sign in to Admin Portal</h2>
            <p className="text-slate-400 text-sm">Enter your credentials to access the dashboard</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <Alert variant="destructive" className="border-red-500/50 bg-red-500/10">
                <AlertDescription className="text-red-400">{error}</AlertDescription>
              </Alert>
            )}

            <div className="space-y-1.5">
              <Label htmlFor="email" className="text-slate-300 text-sm font-medium">Email Address</Label>
              <Input
                id="email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                placeholder="admin@uduxpass.com"
                autoComplete="email"
                className="h-11 text-white placeholder:text-slate-500"
                style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)' }}
              />
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="password" className="text-slate-300 text-sm font-medium">Password</Label>
              <div className="relative">
                <Input
                  id="password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  value={formData.password}
                  onChange={handleChange}
                  placeholder="Enter your password"
                  autoComplete="current-password"
                  className="h-11 pr-10 text-white placeholder:text-slate-500"
                  style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)' }}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200 transition-colors"
                >
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
            </div>

            <Button
              type="submit"
              disabled={isLoading}
              className="w-full h-11 font-semibold text-slate-900 transition-all"
              style={{ background: isLoading ? '#d97706' : '#F59E0B', color: '#0f1729' }}
            >
              {isLoading ? (
                <span className="flex items-center gap-2">
                  <span className="w-4 h-4 border-2 border-slate-900/30 border-t-slate-900 rounded-full animate-spin" />
                  Signing in...
                </span>
              ) : 'Sign In to Admin Portal'}
            </Button>
          </form>

          {/* Demo credentials */}
          <div className="mt-8 p-4 rounded-xl" style={{ background: 'rgba(245,158,11,0.06)', border: '1px solid rgba(245,158,11,0.15)' }}>
            <h4 className="text-amber-400 text-xs font-semibold uppercase tracking-wider mb-3">Demo Credentials</h4>
            <div className="space-y-1.5 text-xs">
              <div className="flex justify-between">
                <span className="text-slate-400">Super Admin</span>
                <span className="text-slate-300 font-mono">admin@uduxpass.com / Admin@123!</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Event Manager</span>
                <span className="text-slate-300 font-mono">eventmanager@uduxpass.com / Admin@123!</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Support Agent</span>
                <span className="text-slate-300 font-mono">support@uduxpass.com / Admin@123!</span>
              </div>
              <div className="flex justify-between">
                <span className="text-slate-400">Data Analyst</span>
                <span className="text-slate-300 font-mono">analyst@uduxpass.com / Admin@123!</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default AdminLoginPage

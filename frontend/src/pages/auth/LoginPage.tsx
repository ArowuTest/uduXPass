import React, { useState } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAuth } from '../../contexts/AuthContext'
import { Eye, EyeOff, Mail, Lock, Smartphone, ArrowLeft } from 'lucide-react'
import { motion } from 'framer-motion'
import { toast } from '@/components/ui/toaster'

// TypeScript interfaces
interface LoginFormData {
  email: string
  password: string
  phone: string
}

interface LoginResult {
  success: boolean
  error?: string
  user?: any
}

interface LocationState {
  from?: {
    pathname: string
  }
}

type LoginMethod = 'email' | 'momo'

const LoginPage: React.FC = () => {
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
    phone: ''
  })
  const [showPassword, setShowPassword] = useState<boolean>(false)
  const [loginMethod, setLoginMethod] = useState<LoginMethod>('email')
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const [otpSent, setOtpSent] = useState<boolean>(false)
  const [otp, setOtp] = useState<string>('')

  const { login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const from = (location.state as LocationState)?.from?.pathname || '/'

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }))
  }

  const handleEmailLogin = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault()
    
    console.log('Email login form submitted', { email: formData.email });
    
    if (!formData.email || !formData.password) {
      console.error('Email login validation failed: missing fields');
      toast.error('Please fill in all fields')
      return
    }

    console.log('Email login validation passed, calling login API...');
    setIsLoading(true)
    
    try {
      await login(formData.email, formData.password)
      // If we get here, login succeeded
      toast.success('Login successful!')
      navigate(from, { replace: true })
    } catch (error) {
      console.error('Login network error:', error);
      toast.error(error instanceof Error ? error.message : 'An error occurred during login')
    } finally {
      setIsLoading(false)
    }
  }

  const handleMomoLogin = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault()
    
    if (!formData.phone) {
      toast.error('Please enter your phone number')
      return
    }

    // Relaxed Nigerian phone validation: accepts +234, 234, or 0 prefix followed by 10 digits
    const phoneRegex = /^(\+234|234|0)\d{10}$/
    const cleanedPhone = formData.phone.replace(/\s+/g, ''); // Remove spaces
    if (!phoneRegex.test(cleanedPhone)) {
      console.error('MoMo phone validation failed:', formData.phone);
      toast.error('Please enter a valid Nigerian phone number (e.g., +2348012345678)')
      return
    }

    setIsLoading(true)
    
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/auth/momo/request-otp`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ phone: formData.phone })
      })

      if (response.ok) {
        setOtpSent(true)
        toast.success('OTP sent to your phone!')
      } else {
        const error = await response.json()
        toast.error(error.message || 'Failed to send OTP')
      }
    } catch (error) {
      toast.error('Failed to send OTP')
    } finally {
      setIsLoading(false)
    }
  }

  const handleOtpVerification = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault()
    
    if (!otp || otp.length !== 6) {
      toast.error('Please enter a valid 6-digit OTP')
      return
    }

    setIsLoading(true)
    
    try {
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/auth/momo/verify-otp`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ 
          phone: formData.phone,
          otp: otp
        })
      })

      if (response.ok) {
        const result = await response.json()
        if (result.success) {
          // Store auth token
          localStorage.setItem('authToken', result.data.token)
          alert('Login successful!')
          navigate(from, { replace: true })
        } else {
          alert(result.error || 'OTP verification failed')
        }
      } else {
        const error = await response.json()
        alert(error.message || 'OTP verification failed')
      }
    } catch (error) {
      alert('An error occurred during OTP verification')
    } finally {
      setIsLoading(false)
    }
  }

  const resetOtpFlow = (): void => {
    setOtpSent(false)
    setOtp('')
    setFormData(prev => ({ ...prev, phone: '' }))
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center p-4">
      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md"
      >
        {/* Back to home */}
        <div className="mb-6">
          <Link 
            to="/" 
            className="inline-flex items-center text-white/70 hover:text-white transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to uduXPass
          </Link>
        </div>

        <Card className="border-slate-700 bg-slate-800/50 backdrop-blur-sm">
          <CardHeader className="text-center">
            <CardTitle className="text-2xl font-bold text-white">Welcome Back</CardTitle>
            <p className="text-slate-300">Sign in to your uduXPass account</p>
          </CardHeader>
          <CardContent>
            <Tabs value={loginMethod} onValueChange={(value: string) => setLoginMethod(value as LoginMethod)}>
              <TabsList className="grid w-full grid-cols-2 mb-6">
                <TabsTrigger value="email" className="flex items-center">
                  <Mail className="w-4 h-4 mr-2" />
                  Email
                </TabsTrigger>
                <TabsTrigger value="momo" className="flex items-center">
                  <Smartphone className="w-4 h-4 mr-2" />
                  MoMo PSB
                </TabsTrigger>
              </TabsList>

              <TabsContent value="email">
                <form onSubmit={handleEmailLogin} className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="email" className="text-slate-200">Email</Label>
                    <Input
                      id="email"
                      name="email"
                      type="email"
                      value={formData.email}
                      onChange={handleInputChange}
                      placeholder="your@email.com"
                      required
                      className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
                    />
                  </div>
                  
                  <div className="space-y-2">
                    <Label htmlFor="password" className="text-slate-200">Password</Label>
                    <div className="relative">
                      <Input
                        id="password"
                        name="password"
                        type={showPassword ? 'text' : 'password'}
                        value={formData.password}
                        onChange={handleInputChange}
                        placeholder="Enter your password"
                        required
                        className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400 pr-10"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="absolute right-0 top-0 h-full px-3 text-slate-400 hover:text-white"
                        onClick={() => setShowPassword(!showPassword)}
                      >
                        {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                      </Button>
                    </div>
                  </div>
                  
                  <Button 
                    type="submit" 
                    className="w-full bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700"
                    disabled={isLoading}
                  >
                    {isLoading ? 'Signing in...' : 'Sign In'}
                  </Button>
                </form>
              </TabsContent>

              <TabsContent value="momo">
                {!otpSent ? (
                  <form onSubmit={handleMomoLogin} className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="phone" className="text-slate-200">Phone Number</Label>
                      <Input
                        id="phone"
                        name="phone"
                        type="tel"
                        value={formData.phone}
                        onChange={handleInputChange}
                        placeholder="+234 801 234 5678"
                        required
                        className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400"
                      />
                      <p className="text-xs text-slate-400">
                        We'll send a verification code to your MoMo PSB registered number
                      </p>
                    </div>
                    
                    <Button 
                      type="submit" 
                      className="w-full bg-gradient-to-r from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700"
                      disabled={isLoading}
                    >
                      {isLoading ? 'Sending OTP...' : 'Send OTP'}
                    </Button>
                  </form>
                ) : (
                  <form onSubmit={handleOtpVerification} className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="otp" className="text-slate-200">Verification Code</Label>
                      <Input
                        id="otp"
                        name="otp"
                        type="text"
                        value={otp}
                        onChange={(e) => setOtp(e.target.value)}
                        placeholder="Enter 6-digit code"
                        maxLength={6}
                        required
                        className="bg-slate-700 border-slate-600 text-white placeholder:text-slate-400 text-center text-lg tracking-widest"
                      />
                      <p className="text-xs text-slate-400">
                        Code sent to {formData.phone}
                      </p>
                    </div>
                    
                    <Button 
                      type="submit" 
                      className="w-full bg-gradient-to-r from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700"
                      disabled={isLoading}
                    >
                      {isLoading ? 'Verifying...' : 'Verify & Sign In'}
                    </Button>
                    
                    <Button 
                      type="button"
                      variant="ghost"
                      className="w-full text-slate-300 hover:text-white"
                      onClick={resetOtpFlow}
                    >
                      Use different number
                    </Button>
                  </form>
                )}
              </TabsContent>
            </Tabs>

            <div className="mt-6 text-center">
              <p className="text-slate-400 text-sm">
                Don't have an account?{' '}
                <Link to="/register" className="text-purple-400 hover:text-purple-300 font-medium">
                  Sign up
                </Link>
              </p>
            </div>

            {/* Demo credentials */}
            <div className="mt-6 p-4 bg-slate-700/50 rounded-lg">
              <h4 className="text-sm font-medium text-slate-200 mb-2">Demo Credentials:</h4>
              <div className="text-xs text-slate-300 space-y-1">
                <div><strong>Email:</strong> user@uduxpass.com / password123</div>
                <div><strong>MoMo:</strong> +234 801 234 5678 / Any 6-digit OTP</div>
              </div>
            </div>
          </CardContent>
        </Card>
      </motion.div>
    </div>
  )
}

export default LoginPage


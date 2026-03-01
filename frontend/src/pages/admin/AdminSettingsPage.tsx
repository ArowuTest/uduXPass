import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Settings,
  Save,
  RefreshCw,
  Globe,
  CreditCard,
  Mail,
  Bell,
  Shield,
  Database,
  Smartphone,
  Users,
  Calendar,
  DollarSign,
  Lock,
  Eye,
  EyeOff,
  AlertTriangle,
  CheckCircle
} from 'lucide-react'

// TypeScript interfaces
interface PaymentSettings {
  momo_psb_enabled: boolean
  momo_psb_api_key: string
  momo_psb_secret_key: string
  paystack_enabled: boolean
  paystack_public_key: string
  paystack_secret_key: string
  default_currency: string
  transaction_fee_percentage: number
}

interface EmailSettings {
  smtp_host: string
  smtp_port: number
  smtp_username: string
  smtp_password: string
  from_email: string
  from_name: string
  email_verification_enabled: boolean
  welcome_email_enabled: boolean
}

interface NotificationSettings {
  push_notifications_enabled: boolean
  email_notifications_enabled: boolean
  sms_notifications_enabled: boolean
  order_confirmation_enabled: boolean
  event_reminders_enabled: boolean
  marketing_emails_enabled: boolean
}

interface SecuritySettings {
  two_factor_enabled: boolean
  password_min_length: number
  password_require_uppercase: boolean
  password_require_lowercase: boolean
  password_require_numbers: boolean
  password_require_symbols: boolean
  session_timeout_minutes: number
  max_login_attempts: number
}

interface PlatformSettings {
  platform_name: string
  platform_description: string
  support_email: string
  support_phone: string
  terms_url: string
  privacy_url: string
  logo_url: string
  favicon_url: string
  maintenance_mode: boolean
  registration_enabled: boolean
}

interface AllSettings {
  payment: PaymentSettings
  email: EmailSettings
  notifications: NotificationSettings
  security: SecuritySettings
  platform: PlatformSettings
}

interface SaveStatus {
  type: 'success' | 'error'
  message: string
}

const AdminSettingsPage: React.FC = () => {
  const [settings, setSettings] = useState<AllSettings | null>(null)
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [isSaving, setIsSaving] = useState<boolean>(false)
  const [showApiKeys, setShowApiKeys] = useState<boolean>(false)
  const [saveStatus, setSaveStatus] = useState<SaveStatus | null>(null)

  useEffect(() => {
    fetchSettings()
  }, [])

  const fetchSettings = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('/v1/admin/settings', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: AllSettings } = await response.json()
        setSettings(result.data || mockSettings)
      } else {
        setSettings(mockSettings)
      }
    } catch (error) {
      console.error('Failed to fetch settings:', error)
      setSettings(mockSettings)
    } finally {
      setIsLoading(false)
    }
  }

  const handleSaveSettings = async (section: keyof AllSettings): Promise<void> => {
    if (!settings) return
    
    setIsSaving(true)
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`/v1/admin/settings/${section}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(settings[section])
      })
      
      if (response.ok) {
        setSaveStatus({ type: 'success', message: 'Settings saved successfully' })
      } else {
        setSaveStatus({ type: 'error', message: 'Failed to save settings' })
      }
    } catch (error) {
      console.error('Failed to save settings:', error)
      setSaveStatus({ type: 'error', message: 'An error occurred while saving' })
    } finally {
      setIsSaving(false)
      setTimeout(() => setSaveStatus(null), 3000)
    }
  }

  const updateSetting = <T extends keyof AllSettings, K extends keyof AllSettings[T]>(
    section: T, 
    key: K, 
    value: AllSettings[T][K]
  ): void => {
    if (!settings) return
    
    setSettings(prev => ({
      ...prev!,
      [section]: {
        ...prev![section],
        [key]: value
      }
    }))
  }

  const testEmailConnection = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('/v1/admin/settings/email/test', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        alert('Email connection test successful!')
      } else {
        alert('Email connection test failed')
      }
    } catch (error) {
      alert('Failed to test email connection')
    }
  }

  // Mock settings for demonstration
  const mockSettings: AllSettings = {
    payment: {
      momo_psb_enabled: true,
      momo_psb_api_key: 'momo_test_key_123',
      momo_psb_secret_key: 'momo_test_secret_456',
      paystack_enabled: true,
      paystack_public_key: 'pk_test_paystack_123',
      paystack_secret_key: 'sk_test_paystack_456',
      default_currency: 'NGN',
      transaction_fee_percentage: 2.5
    },
    email: {
      smtp_host: 'smtp.gmail.com',
      smtp_port: 587,
      smtp_username: 'noreply@uduxpass.com',
      smtp_password: 'email_password_123',
      from_email: 'noreply@uduxpass.com',
      from_name: 'uduXPass',
      email_verification_enabled: true,
      welcome_email_enabled: true
    },
    notifications: {
      push_notifications_enabled: true,
      email_notifications_enabled: true,
      sms_notifications_enabled: false,
      order_confirmation_enabled: true,
      event_reminders_enabled: true,
      marketing_emails_enabled: false
    },
    security: {
      two_factor_enabled: false,
      password_min_length: 8,
      password_require_uppercase: true,
      password_require_lowercase: true,
      password_require_numbers: true,
      password_require_symbols: false,
      session_timeout_minutes: 60,
      max_login_attempts: 5
    },
    platform: {
      platform_name: 'uduXPass',
      platform_description: 'Premium ticketing platform for exclusive events',
      support_email: 'support@uduxpass.com',
      support_phone: '+234 800 123 4567',
      terms_url: 'https://uduxpass.com/terms',
      privacy_url: 'https://uduxpass.com/privacy',
      logo_url: '/assets/logo.png',
      favicon_url: '/assets/favicon.ico',
      maintenance_mode: false,
      registration_enabled: true
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-2">Loading settings...</p>
        </div>
      </div>
    )
  }

  if (!settings) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <p className="">Failed to load settings</p>
          <Button onClick={fetchSettings} className="mt-4">
            <RefreshCw className="h-4 w-4 mr-2" />
            Retry
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
          <p className="text-muted-foreground">
            Configure platform settings and integrations
          </p>
        </div>
        <Button variant="outline" onClick={fetchSettings}>
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      {/* Save Status Alert */}
      {saveStatus && (
        <Alert variant={saveStatus.type === 'error' ? 'destructive' : 'default'}>
          {saveStatus.type === 'success' ? (
            <CheckCircle className="h-4 w-4" />
          ) : (
            <AlertTriangle className="h-4 w-4" />
          )}
          <AlertDescription>{saveStatus.message}</AlertDescription>
        </Alert>
      )}

      <Tabs defaultValue="platform" className="space-y-4">
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="platform">Platform</TabsTrigger>
          <TabsTrigger value="payment">Payment</TabsTrigger>
          <TabsTrigger value="email">Email</TabsTrigger>
          <TabsTrigger value="notifications">Notifications</TabsTrigger>
          <TabsTrigger value="security">Security</TabsTrigger>
        </TabsList>

        {/* Platform Settings */}
        <TabsContent value="platform" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Globe className="h-5 w-5 mr-2" />
                Platform Configuration
              </CardTitle>
              <CardDescription>
                Basic platform information and branding
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="platform_name">Platform Name</Label>
                  <Input
                    id="platform_name"
                    value={settings.platform.platform_name}
                    onChange={(e) => updateSetting('platform', 'platform_name', e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="support_email">Support Email</Label>
                  <Input
                    id="support_email"
                    type="email"
                    value={settings.platform.support_email}
                    onChange={(e) => updateSetting('platform', 'support_email', e.target.value)}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="platform_description">Platform Description</Label>
                <Textarea
                  id="platform_description"
                  value={settings.platform.platform_description}
                  onChange={(e) => updateSetting('platform', 'platform_description', e.target.value)}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="support_phone">Support Phone</Label>
                  <Input
                    id="support_phone"
                    value={settings.platform.support_phone}
                    onChange={(e) => updateSetting('platform', 'support_phone', e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="logo_url">Logo URL</Label>
                  <Input
                    id="logo_url"
                    value={settings.platform.logo_url}
                    onChange={(e) => updateSetting('platform', 'logo_url', e.target.value)}
                  />
                </div>
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Maintenance Mode</Label>
                  <p className="text-sm text-muted-foreground">
                    Temporarily disable public access to the platform
                  </p>
                </div>
                <Switch
                  checked={settings.platform.maintenance_mode}
                  onCheckedChange={(checked) => updateSetting('platform', 'maintenance_mode', checked)}
                />
              </div>

              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>User Registration</Label>
                  <p className="text-sm text-muted-foreground">
                    Allow new users to register accounts
                  </p>
                </div>
                <Switch
                  checked={settings.platform.registration_enabled}
                  onCheckedChange={(checked) => updateSetting('platform', 'registration_enabled', checked)}
                />
              </div>

              <Button onClick={() => handleSaveSettings('platform')} disabled={isSaving}>
                <Save className="h-4 w-4 mr-2" />
                {isSaving ? 'Saving...' : 'Save Platform Settings'}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Payment Settings */}
        <TabsContent value="payment" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <CreditCard className="h-5 w-5 mr-2" />
                Payment Gateway Configuration
              </CardTitle>
              <CardDescription>
                Configure MoMo PSB and Paystack payment integrations
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* MoMo PSB Settings */}
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium">MoMo PSB Integration</h4>
                    <p className="text-sm text-muted-foreground">
                      Mobile Money Payment Service Bank
                    </p>
                  </div>
                  <Switch
                    checked={settings.payment.momo_psb_enabled}
                    onCheckedChange={(checked) => updateSetting('payment', 'momo_psb_enabled', checked)}
                  />
                </div>

                {settings.payment.momo_psb_enabled && (
                  <div className="grid grid-cols-2 gap-4 pl-4 border-l-2 border-blue-200">
                    <div className="space-y-2">
                      <Label htmlFor="momo_api_key">API Key</Label>
                      <div className="relative">
                        <Input
                          id="momo_api_key"
                          type={showApiKeys ? 'text' : 'password'}
                          value={settings.payment.momo_psb_api_key}
                          onChange={(e) => updateSetting('payment', 'momo_psb_api_key', e.target.value)}
                        />
                      </div>
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="momo_secret_key">Secret Key</Label>
                      <Input
                        id="momo_secret_key"
                        type={showApiKeys ? 'text' : 'password'}
                        value={settings.payment.momo_psb_secret_key}
                        onChange={(e) => updateSetting('payment', 'momo_psb_secret_key', e.target.value)}
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Paystack Settings */}
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium">Paystack Integration</h4>
                    <p className="text-sm text-muted-foreground">
                      Alternative payment gateway
                    </p>
                  </div>
                  <Switch
                    checked={settings.payment.paystack_enabled}
                    onCheckedChange={(checked) => updateSetting('payment', 'paystack_enabled', checked)}
                  />
                </div>

                {settings.payment.paystack_enabled && (
                  <div className="grid grid-cols-2 gap-4 pl-4 border-l-2 border-green-200">
                    <div className="space-y-2">
                      <Label htmlFor="paystack_public_key">Public Key</Label>
                      <Input
                        id="paystack_public_key"
                        value={settings.payment.paystack_public_key}
                        onChange={(e) => updateSetting('payment', 'paystack_public_key', e.target.value)}
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="paystack_secret_key">Secret Key</Label>
                      <Input
                        id="paystack_secret_key"
                        type={showApiKeys ? 'text' : 'password'}
                        value={settings.payment.paystack_secret_key}
                        onChange={(e) => updateSetting('payment', 'paystack_secret_key', e.target.value)}
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* General Payment Settings */}
              <div className="space-y-4">
                <h4 className="font-medium">General Settings</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="default_currency">Default Currency</Label>
                    <Select
                      value={settings.payment.default_currency}
                      onValueChange={(value) => updateSetting('payment', 'default_currency', value)}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="NGN">Nigerian Naira (NGN)</SelectItem>
                        <SelectItem value="USD">US Dollar (USD)</SelectItem>
                        <SelectItem value="EUR">Euro (EUR)</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="transaction_fee">Transaction Fee (%)</Label>
                    <Input
                      id="transaction_fee"
                      type="number"
                      step="0.1"
                      value={settings.payment.transaction_fee_percentage}
                      onChange={(e) => updateSetting('payment', 'transaction_fee_percentage', parseFloat(e.target.value))}
                    />
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <Button
                  variant="outline"
                  onClick={() => setShowApiKeys(!showApiKeys)}
                >
                  {showApiKeys ? <EyeOff className="h-4 w-4 mr-2" /> : <Eye className="h-4 w-4 mr-2" />}
                  {showApiKeys ? 'Hide' : 'Show'} API Keys
                </Button>
                <Button onClick={() => handleSaveSettings('payment')} disabled={isSaving}>
                  <Save className="h-4 w-4 mr-2" />
                  {isSaving ? 'Saving...' : 'Save Payment Settings'}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Email Settings */}
        <TabsContent value="email" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Mail className="h-5 w-5 mr-2" />
                Email Configuration
              </CardTitle>
              <CardDescription>
                Configure SMTP settings for email delivery
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="smtp_host">SMTP Host</Label>
                  <Input
                    id="smtp_host"
                    value={settings.email.smtp_host}
                    onChange={(e) => updateSetting('email', 'smtp_host', e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="smtp_port">SMTP Port</Label>
                  <Input
                    id="smtp_port"
                    type="number"
                    value={settings.email.smtp_port}
                    onChange={(e) => updateSetting('email', 'smtp_port', parseInt(e.target.value))}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="smtp_username">SMTP Username</Label>
                  <Input
                    id="smtp_username"
                    value={settings.email.smtp_username}
                    onChange={(e) => updateSetting('email', 'smtp_username', e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="smtp_password">SMTP Password</Label>
                  <Input
                    id="smtp_password"
                    type="password"
                    value={settings.email.smtp_password}
                    onChange={(e) => updateSetting('email', 'smtp_password', e.target.value)}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="from_email">From Email</Label>
                  <Input
                    id="from_email"
                    type="email"
                    value={settings.email.from_email}
                    onChange={(e) => updateSetting('email', 'from_email', e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="from_name">From Name</Label>
                  <Input
                    id="from_name"
                    value={settings.email.from_name}
                    onChange={(e) => updateSetting('email', 'from_name', e.target.value)}
                  />
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Email Verification</Label>
                    <p className="text-sm text-muted-foreground">
                      Require email verification for new accounts
                    </p>
                  </div>
                  <Switch
                    checked={settings.email.email_verification_enabled}
                    onCheckedChange={(checked) => updateSetting('email', 'email_verification_enabled', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Welcome Emails</Label>
                    <p className="text-sm text-muted-foreground">
                      Send welcome emails to new users
                    </p>
                  </div>
                  <Switch
                    checked={settings.email.welcome_email_enabled}
                    onCheckedChange={(checked) => updateSetting('email', 'welcome_email_enabled', checked)}
                  />
                </div>
              </div>

              <div className="flex items-center space-x-2">
                <Button variant="outline" onClick={testEmailConnection}>
                  <Mail className="h-4 w-4 mr-2" />
                  Test Connection
                </Button>
                <Button onClick={() => handleSaveSettings('email')} disabled={isSaving}>
                  <Save className="h-4 w-4 mr-2" />
                  {isSaving ? 'Saving...' : 'Save Email Settings'}
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Notifications Settings */}
        <TabsContent value="notifications" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Bell className="h-5 w-5 mr-2" />
                Notification Settings
              </CardTitle>
              <CardDescription>
                Configure notification preferences and channels
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Push Notifications</Label>
                    <p className="text-sm text-muted-foreground">
                      Enable browser push notifications
                    </p>
                  </div>
                  <Switch
                    checked={settings.notifications.push_notifications_enabled}
                    onCheckedChange={(checked) => updateSetting('notifications', 'push_notifications_enabled', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Email Notifications</Label>
                    <p className="text-sm text-muted-foreground">
                      Send notifications via email
                    </p>
                  </div>
                  <Switch
                    checked={settings.notifications.email_notifications_enabled}
                    onCheckedChange={(checked) => updateSetting('notifications', 'email_notifications_enabled', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>SMS Notifications</Label>
                    <p className="text-sm text-muted-foreground">
                      Send notifications via SMS
                    </p>
                  </div>
                  <Switch
                    checked={settings.notifications.sms_notifications_enabled}
                    onCheckedChange={(checked) => updateSetting('notifications', 'sms_notifications_enabled', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Order Confirmations</Label>
                    <p className="text-sm text-muted-foreground">
                      Send order confirmation notifications
                    </p>
                  </div>
                  <Switch
                    checked={settings.notifications.order_confirmation_enabled}
                    onCheckedChange={(checked) => updateSetting('notifications', 'order_confirmation_enabled', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Event Reminders</Label>
                    <p className="text-sm text-muted-foreground">
                      Send event reminder notifications
                    </p>
                  </div>
                  <Switch
                    checked={settings.notifications.event_reminders_enabled}
                    onCheckedChange={(checked) => updateSetting('notifications', 'event_reminders_enabled', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <Label>Marketing Emails</Label>
                    <p className="text-sm text-muted-foreground">
                      Send promotional and marketing emails
                    </p>
                  </div>
                  <Switch
                    checked={settings.notifications.marketing_emails_enabled}
                    onCheckedChange={(checked) => updateSetting('notifications', 'marketing_emails_enabled', checked)}
                  />
                </div>
              </div>

              <Button onClick={() => handleSaveSettings('notifications')} disabled={isSaving}>
                <Save className="h-4 w-4 mr-2" />
                {isSaving ? 'Saving...' : 'Save Notification Settings'}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Security Settings */}
        <TabsContent value="security" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Shield className="h-5 w-5 mr-2" />
                Security Configuration
              </CardTitle>
              <CardDescription>
                Configure security policies and authentication requirements
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label>Two-Factor Authentication</Label>
                  <p className="text-sm text-muted-foreground">
                    Require 2FA for admin accounts
                  </p>
                </div>
                <Switch
                  checked={settings.security.two_factor_enabled}
                  onCheckedChange={(checked) => updateSetting('security', 'two_factor_enabled', checked)}
                />
              </div>

              <div className="space-y-4">
                <h4 className="font-medium">Password Requirements</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="password_min_length">Minimum Length</Label>
                    <Input
                      id="password_min_length"
                      type="number"
                      min="6"
                      max="32"
                      value={settings.security.password_min_length}
                      onChange={(e) => updateSetting('security', 'password_min_length', parseInt(e.target.value))}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="session_timeout">Session Timeout (minutes)</Label>
                    <Input
                      id="session_timeout"
                      type="number"
                      min="15"
                      max="480"
                      value={settings.security.session_timeout_minutes}
                      onChange={(e) => updateSetting('security', 'session_timeout_minutes', parseInt(e.target.value))}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <Label>Require Uppercase Letters</Label>
                    <Switch
                      checked={settings.security.password_require_uppercase}
                      onCheckedChange={(checked) => updateSetting('security', 'password_require_uppercase', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <Label>Require Lowercase Letters</Label>
                    <Switch
                      checked={settings.security.password_require_lowercase}
                      onCheckedChange={(checked) => updateSetting('security', 'password_require_lowercase', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <Label>Require Numbers</Label>
                    <Switch
                      checked={settings.security.password_require_numbers}
                      onCheckedChange={(checked) => updateSetting('security', 'password_require_numbers', checked)}
                    />
                  </div>
                  <div className="flex items-center justify-between">
                    <Label>Require Symbols</Label>
                    <Switch
                      checked={settings.security.password_require_symbols}
                      onCheckedChange={(checked) => updateSetting('security', 'password_require_symbols', checked)}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="max_login_attempts">Max Login Attempts</Label>
                  <Input
                    id="max_login_attempts"
                    type="number"
                    min="3"
                    max="10"
                    value={settings.security.max_login_attempts}
                    onChange={(e) => updateSetting('security', 'max_login_attempts', parseInt(e.target.value))}
                  />
                  <p className="text-sm text-muted-foreground">
                    Account will be locked after this many failed attempts
                  </p>
                </div>
              </div>

              <Button onClick={() => handleSaveSettings('security')} disabled={isSaving}>
                <Save className="h-4 w-4 mr-2" />
                {isSaving ? 'Saving...' : 'Save Security Settings'}
              </Button>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default AdminSettingsPage


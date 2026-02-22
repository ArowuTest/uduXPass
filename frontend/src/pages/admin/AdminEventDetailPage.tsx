import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  ArrowLeft,
  Calendar, 
  MapPin, 
  Users, 
  DollarSign,
  CreditCard,
  Smartphone,
  Building2,
  Settings,
  Save,
  Eye,
  Edit,
  CheckCircle,
  Globe
} from 'lucide-react'

interface Event {
  id: string
  name: string
  description: string
  startDate: string
  venue: {
    name: string
    city: string
    address: string
  }
  status: string
  ticketsSold: number
  revenue: number
  capacity: number
  ticketTiers: Array<{
    name: string
    price: number
    quantity: number
    sold: number
  }>
}

interface PaymentSettings {
  cardPayments: boolean
  bankTransfer: boolean
  mobileMoney: boolean
  paystack: {
    enabled: boolean
    publicKey: string
    secretKey: string
  }
  flutterwave: {
    enabled: boolean
    publicKey: string
    secretKey: string
  }
  monnify: {
    enabled: boolean
    apiKey: string
    secretKey: string
    contractCode: string
  }
}

const AdminEventDetailPageEnhanced: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  
  const [event, setEvent] = useState<Event | null>(null)
  const [paymentSettings, setPaymentSettings] = useState<PaymentSettings>({
    cardPayments: true,
    bankTransfer: false,
    mobileMoney: true,
    paystack: {
      enabled: true,
      publicKey: 'pk_test_xxxxxxxxxxxxx',
      secretKey: 'sk_test_xxxxxxxxxxxxx'
    },
    flutterwave: {
      enabled: false,
      publicKey: '',
      secretKey: ''
    },
    monnify: {
      enabled: false,
      apiKey: '',
      secretKey: '',
      contractCode: ''
    }
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [publishing, setPublishing] = useState(false)

  // Mock event data - in real app this would come from API
  const mockEvent: Event = {
    id: 'event-3',
    name: 'Wizkid Live in Lagos',
    description: 'The Lagos grand finale of the Wizkid tour featuring special guest appearances and exclusive performances. Experience the biggest Afrobeats concert of the year at the iconic Eko Atlantic City.',
    startDate: '2025-10-18T20:00:00Z',
    venue: {
      name: 'Eko Atlantic City',
      city: 'Lagos',
      address: 'Eko Atlantic City, Victoria Island, Lagos'
    },
    status: 'upcoming',
    ticketsSold: 280,
    revenue: 3200000,
    capacity: 1700,
    ticketTiers: [
      {
        name: 'VIP Diamond',
        price: 75000,
        quantity: 200,
        sold: 45
      },
      {
        name: 'Golden Circle',
        price: 30000,
        quantity: 1500,
        sold: 235
      }
    ]
  }

  useEffect(() => {
    // Simulate API call
    setTimeout(() => {
      setEvent(mockEvent)
      setLoading(false)
    }, 1000)
  }, [id])

  const handlePaymentSettingsChange = (key: string, value: any) => {
    setPaymentSettings(prev => ({
      ...prev,
      [key]: value
    }))
  }

  const handleProviderSettingsChange = (provider: string, key: string, value: any) => {
    setPaymentSettings(prev => ({
      ...prev,
      [provider]: {
        ...prev[provider as keyof PaymentSettings] as any,
        [key]: value
      }
    }))
  }

  const handleSavePaymentSettings = async () => {
    setSaving(true)
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    alert('✅ Payment settings saved successfully! Event is now configured for ticket sales.')
    setSaving(false)
  }

  const handlePublishEvent = async () => {
    setPublishing(true)
    
    // Simulate API call to publish event
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Update event status
    if (event) {
      setEvent({
        ...event,
        status: 'published'
      })
    }
    
    alert('🎉 Event published successfully! Customers can now purchase tickets.')
    setPublishing(false)
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: 'NGN',
      minimumFractionDigits: 0
    }).format(amount)
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (!event) {
    return (
      <div className="text-center py-12">
        <h3 className="text-lg font-medium mb-2">Event not found</h3>
        <Button onClick={() => navigate('/admin/events')}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Events
        </Button>
      </div>
    )
  }

  const getStatusBadge = (status: string) => {
    const statusConfig = {
      published: { label: 'Published', variant: 'default' as const, color: 'bg-green-100 text-green-800' },
      upcoming: { label: 'Upcoming', variant: 'outline' as const, color: 'bg-yellow-100 text-yellow-800' },
      draft: { label: 'Draft', variant: 'secondary' as const, color: 'bg-gray-100 text-gray-800' }
    }
    
    const config = statusConfig[status as keyof typeof statusConfig] || { label: status, variant: 'outline' as const, color: 'bg-gray-100 text-gray-800' }
    return (
      <div className={`px-3 py-1 rounded-full text-sm font-medium ${config.color}`}>
        {config.label}
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button variant="outline" onClick={() => navigate('/admin/events')}>
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Events
          </Button>
          <div>
            <h1 className="text-3xl font-bold">{event.name}</h1>
            <p className="text-muted-foreground">Event management and configuration</p>
          </div>
        </div>
        <div className="flex items-center space-x-3">
          {getStatusBadge(event.status)}
          {event.status !== 'published' && (
            <Button 
              onClick={handlePublishEvent} 
              disabled={publishing}
              className="bg-green-600 hover:bg-green-700"
            >
              <Globe className="h-4 w-4 mr-2" />
              {publishing ? 'Publishing...' : 'Publish Event'}
            </Button>
          )}
        </div>
      </div>

      {/* Event Overview */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Tickets Sold</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{event.ticketsSold}</div>
            <p className="text-xs text-muted-foreground">of {event.capacity} capacity</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(event.revenue)}</div>
            <p className="text-xs text-muted-foreground">Total sales</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Event Date</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-sm font-bold">{formatDate(event.startDate)}</div>
            <p className="text-xs text-muted-foreground">Start time</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Venue</CardTitle>
            <MapPin className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-sm font-bold">{event.venue.name}</div>
            <p className="text-xs text-muted-foreground">{event.venue.city}</p>
          </CardContent>
        </Card>
      </div>

      {/* Tabs for different sections */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="payments">Payment Settings</TabsTrigger>
          <TabsTrigger value="tickets">Ticket Tiers</TabsTrigger>
          <TabsTrigger value="settings">Event Settings</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Event Details</CardTitle>
              <CardDescription>Basic information about the event</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h4 className="font-medium mb-2">Description</h4>
                <p className="text-sm text-muted-foreground">{event.description}</p>
              </div>
              <div>
                <h4 className="font-medium mb-2">Venue Details</h4>
                <p className="text-sm text-muted-foreground">{event.venue.address}</p>
              </div>
              <div>
                <h4 className="font-medium mb-2">Event Status</h4>
                <div className="flex items-center space-x-2">
                  {getStatusBadge(event.status)}
                  {event.status === 'published' && (
                    <span className="text-sm text-green-600">✅ Live and accepting ticket purchases</span>
                  )}
                  {event.status === 'upcoming' && (
                    <span className="text-sm text-yellow-600">⏳ Ready to publish</span>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="payments" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <CreditCard className="h-5 w-5 mr-2" />
                Payment Methods
              </CardTitle>
              <CardDescription>Configure which payment methods are available for this event</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Payment Method Toggles */}
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <CreditCard className="h-4 w-4" />
                    <div>
                      <p className="font-medium">Card Payments</p>
                      <p className="text-sm text-muted-foreground">Accept Visa, Mastercard, Verve</p>
                    </div>
                  </div>
                  <Switch
                    checked={paymentSettings.cardPayments}
                    onCheckedChange={(checked) => handlePaymentSettingsChange('cardPayments', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Building2 className="h-4 w-4" />
                    <div>
                      <p className="font-medium">Bank Transfer</p>
                      <p className="text-sm text-muted-foreground">Direct bank transfers</p>
                    </div>
                  </div>
                  <Switch
                    checked={paymentSettings.bankTransfer}
                    onCheckedChange={(checked) => handlePaymentSettingsChange('bankTransfer', checked)}
                  />
                </div>

                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Smartphone className="h-4 w-4" />
                    <div>
                      <p className="font-medium">Mobile Money</p>
                      <p className="text-sm text-muted-foreground">MTN, Airtel, 9mobile wallets</p>
                    </div>
                  </div>
                  <Switch
                    checked={paymentSettings.mobileMoney}
                    onCheckedChange={(checked) => handlePaymentSettingsChange('mobileMoney', checked)}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Payment Providers */}
          <Card>
            <CardHeader>
              <CardTitle>Payment Providers</CardTitle>
              <CardDescription>Configure payment gateway settings</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* Paystack */}
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Paystack</p>
                    <p className="text-sm text-muted-foreground">Primary payment processor</p>
                  </div>
                  <Switch
                    checked={paymentSettings.paystack.enabled}
                    onCheckedChange={(checked) => handleProviderSettingsChange('paystack', 'enabled', checked)}
                  />
                </div>
                
                {paymentSettings.paystack.enabled && (
                  <div className="grid gap-4 md:grid-cols-2 pl-4">
                    <div>
                      <label className="text-sm font-medium">Public Key</label>
                      <Input
                        value={paymentSettings.paystack.publicKey}
                        onChange={(e) => handleProviderSettingsChange('paystack', 'publicKey', e.target.value)}
                        placeholder="pk_test_xxxxxxxxxxxxx"
                      />
                    </div>
                    <div>
                      <label className="text-sm font-medium">Secret Key</label>
                      <Input
                        type="password"
                        value={paymentSettings.paystack.secretKey}
                        onChange={(e) => handleProviderSettingsChange('paystack', 'secretKey', e.target.value)}
                        placeholder="sk_test_xxxxxxxxxxxxx"
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Save Button - PROMINENTLY DISPLAYED */}
              <div className="border-t pt-6">
                <div className="flex justify-between items-center">
                  <div>
                    <h4 className="font-medium">Save Configuration</h4>
                    <p className="text-sm text-muted-foreground">Apply payment settings to this event</p>
                  </div>
                  <Button 
                    onClick={handleSavePaymentSettings} 
                    disabled={saving}
                    size="lg"
                    className="bg-blue-600 hover:bg-blue-700"
                  >
                    <Save className="h-4 w-4 mr-2" />
                    {saving ? 'Saving...' : 'Save Payment Settings'}
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="tickets" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Ticket Tiers</CardTitle>
              <CardDescription>Manage ticket pricing and availability</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {event.ticketTiers.map((tier, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                    <div>
                      <h4 className="font-medium">{tier.name}</h4>
                      <p className="text-sm text-muted-foreground">
                        {formatCurrency(tier.price)} • {tier.sold} of {tier.quantity} sold
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">{formatCurrency(tier.price * tier.sold)}</p>
                      <p className="text-sm text-muted-foreground">Revenue</p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="settings" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Event Settings</CardTitle>
              <CardDescription>Additional event configuration options</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Allow Refunds</p>
                    <p className="text-sm text-muted-foreground">Enable ticket refunds for this event</p>
                  </div>
                  <Switch defaultChecked />
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Show Remaining Tickets</p>
                    <p className="text-sm text-muted-foreground">Display ticket availability to customers</p>
                  </div>
                  <Switch defaultChecked />
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">Require Approval</p>
                    <p className="text-sm text-muted-foreground">Manually approve ticket purchases</p>
                  </div>
                  <Switch />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default AdminEventDetailPageEnhanced

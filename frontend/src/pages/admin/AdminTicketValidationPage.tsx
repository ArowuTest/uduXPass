import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  QrCode,
  Scan,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Search,
  Filter,
  Eye,
  Download,
  RefreshCw,
  Clock,
  User,
  Calendar,
  MapPin,
  Ticket,
  Activity,
  BarChart3
} from 'lucide-react'

// TypeScript interfaces
interface TicketData {
  id: string
  qrCode: string
  customerName: string
  customerEmail: string
  eventName: string
  eventDate: string
  venue: string
  ticketTier: string
  price: number
  status: 'valid' | 'used' | 'expired' | 'cancelled'
  purchaseDate: string
  scanHistory: ScanRecord[]
  orderNumber: string
}

interface ScanRecord {
  id: string
  scannedAt: string
  scannedBy: string
  scannerLocation: string
  scannerDevice: string
  result: 'success' | 'failed' | 'duplicate'
}

interface ScanResult {
  success: boolean
  message: string
  ticket?: TicketData
  error?: string
}

interface ValidationStats {
  totalTickets: number
  validTickets: number
  usedTickets: number
  expiredTickets: number
  cancelledTickets: number
  todayScans: number
}

const AdminTicketValidationPage: React.FC = () => {
  const [tickets, setTickets] = useState<TicketData[]>([])
  const [filteredTickets, setFilteredTickets] = useState<TicketData[]>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [searchTerm, setSearchTerm] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [selectedTicket, setSelectedTicket] = useState<TicketData | null>(null)
  const [showTicketDialog, setShowTicketDialog] = useState<boolean>(false)
  const [manualScanCode, setManualScanCode] = useState<string>('')
  const [scanResult, setScanResult] = useState<ScanResult | null>(null)
  const [stats, setStats] = useState<ValidationStats | null>(null)

  useEffect(() => {
    fetchTickets()
    fetchStats()
  }, [])

  useEffect(() => {
    filterTickets()
  }, [tickets, searchTerm, statusFilter])

  const fetchTickets = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/tickets', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: { tickets: TicketData[] } } = await response.json()
        setTickets(result.data.tickets || mockTickets)
      } else {
        setTickets(mockTickets)
      }
    } catch (error) {
      console.error('Failed to fetch tickets:', error)
      setTickets(mockTickets)
    } finally {
      setIsLoading(false)
    }
  }

  const fetchStats = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/tickets/stats', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: ValidationStats } = await response.json()
        setStats(result.data)
      } else {
        setStats(mockStats)
      }
    } catch (error) {
      console.error('Failed to fetch stats:', error)
      setStats(mockStats)
    }
  }

  const filterTickets = (): void => {
    let filtered = tickets

    if (searchTerm) {
      filtered = filtered.filter(ticket => 
        ticket.qrCode.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ticket.customerName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ticket.eventName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        ticket.orderNumber.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter(ticket => ticket.status === statusFilter)
    }

    setFilteredTickets(filtered)
  }

  const handleViewTicket = async (ticketId: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`http://localhost:8080/v1/admin/tickets/${ticketId}`, {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: TicketData } = await response.json()
        setSelectedTicket(result.data)
      } else {
        // Use mock data if API fails
        const mockTicket = mockTickets.find(t => t.id === ticketId)
        setSelectedTicket(mockTicket || null)
      }
      setShowTicketDialog(true)
    } catch (error) {
      console.error('Failed to fetch ticket details:', error)
      const mockTicket = mockTickets.find(t => t.id === ticketId)
      setSelectedTicket(mockTicket || null)
      setShowTicketDialog(true)
    }
  }

  const handleManualScan = async (): Promise<void> => {
    if (!manualScanCode.trim()) {
      setScanResult({
        success: false,
        message: 'Please enter a QR code to scan',
        error: 'Empty QR code'
      })
      return
    }

    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/tickets/validate', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          qrCode: manualScanCode,
          scannerLocation: 'Admin Portal',
          scannerDevice: 'Manual Entry'
        })
      })
      
      if (response.ok) {
        const result: { data: ScanResult } = await response.json()
        setScanResult(result.data)
        
        if (result.data.success) {
          // Refresh tickets to show updated status
          fetchTickets()
          fetchStats()
        }
      } else {
        // Mock validation for demonstration
        const ticket = mockTickets.find(t => t.qrCode === manualScanCode)
        if (ticket) {
          if (ticket.status === 'valid') {
            setScanResult({
              success: true,
              message: 'Ticket validated successfully!',
              ticket: ticket
            })
          } else if (ticket.status === 'used') {
            setScanResult({
              success: false,
              message: 'Ticket has already been used',
              ticket: ticket,
              error: 'Already scanned'
            })
          } else {
            setScanResult({
              success: false,
              message: `Ticket is ${ticket.status}`,
              ticket: ticket,
              error: `Invalid status: ${ticket.status}`
            })
          }
        } else {
          setScanResult({
            success: false,
            message: 'Invalid QR code - ticket not found',
            error: 'Ticket not found'
          })
        }
      }
    } catch (error) {
      console.error('Failed to validate ticket:', error)
      setScanResult({
        success: false,
        message: 'Failed to validate ticket',
        error: 'Network error'
      })
    }
  }

  const handleInvalidateTicket = async (ticketId: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`http://localhost:8080/v1/admin/tickets/${ticketId}/invalidate`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        alert('Ticket invalidated successfully')
        fetchTickets()
        fetchStats()
        setShowTicketDialog(false)
      } else {
        alert('Failed to invalidate ticket')
      }
    } catch (error) {
      console.error('Failed to invalidate ticket:', error)
      alert('An error occurred while invalidating the ticket')
    }
  }

  const exportTicketData = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/tickets/export', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `tickets-export-${new Date().toISOString().split('T')[0]}.csv`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
      } else {
        alert('Failed to export ticket data')
      }
    } catch (error) {
      console.error('Failed to export tickets:', error)
      alert('An error occurred while exporting tickets')
    }
  }

  const getStatusBadge = (status: string): JSX.Element => {
    const statusConfig = {
      valid: { color: 'bg-green-100 text-green-800', icon: CheckCircle },
      used: { color: 'bg-blue-100 text-blue-800', icon: Scan },
      expired: { color: 'bg-red-100 text-red-800', icon: Clock },
      cancelled: { color: 'bg-gray-100 text-gray-800', icon: XCircle }
    }
    
    const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.valid
    const Icon = config.icon
    
    return (
      <Badge className={config.color}>
        <Icon className="w-3 h-3 mr-1" />
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </Badge>
    )
  }

  // Mock data for demonstration
  const mockTickets: TicketData[] = [
    {
      id: 'ticket-1',
      qrCode: 'QR123456789',
      customerName: 'John Doe',
      customerEmail: 'john@example.com',
      eventName: 'Lagos Music Festival',
      eventDate: '2025-12-15T18:00:00Z',
      venue: 'Tafawa Balewa Square',
      ticketTier: 'VIP',
      price: 15000,
      status: 'valid',
      purchaseDate: '2025-09-01T10:00:00Z',
      orderNumber: 'ORD-001',
      scanHistory: []
    },
    {
      id: 'ticket-2',
      qrCode: 'QR987654321',
      customerName: 'Jane Smith',
      customerEmail: 'jane@example.com',
      eventName: 'Lagos Music Festival',
      eventDate: '2025-12-15T18:00:00Z',
      venue: 'Tafawa Balewa Square',
      ticketTier: 'Regular',
      price: 5000,
      status: 'used',
      purchaseDate: '2025-09-02T14:30:00Z',
      orderNumber: 'ORD-002',
      scanHistory: [
        {
          id: 'scan-1',
          scannedAt: '2025-12-15T17:45:00Z',
          scannedBy: 'Scanner Operator 1',
          scannerLocation: 'Main Entrance',
          scannerDevice: 'Scanner-001',
          result: 'success'
        }
      ]
    }
  ]

  const mockStats: ValidationStats = {
    totalTickets: 450,
    validTickets: 320,
    usedTickets: 85,
    expiredTickets: 30,
    cancelledTickets: 15,
    todayScans: 127
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-2">Loading tickets...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Ticket Validation</h1>
          <p className="text-muted-foreground">
            Validate tickets and monitor scanning activity
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" onClick={exportTicketData}>
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button variant="outline" onClick={fetchTickets}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center">
                <Ticket className="h-4 w-4 text-muted-foreground" />
                <div className="ml-2">
                  <p className="text-sm font-medium">Total Tickets</p>
                  <p className="text-2xl font-bold">{stats.totalTickets}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center">
                <CheckCircle className="h-4 w-4 text-green-600" />
                <div className="ml-2">
                  <p className="text-sm font-medium">Valid</p>
                  <p className="text-2xl font-bold">{stats.validTickets}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center">
                <Scan className="h-4 w-4 text-blue-600" />
                <div className="ml-2">
                  <p className="text-sm font-medium">Used</p>
                  <p className="text-2xl font-bold">{stats.usedTickets}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center">
                <Clock className="h-4 w-4 text-red-600" />
                <div className="ml-2">
                  <p className="text-sm font-medium">Expired</p>
                  <p className="text-2xl font-bold">{stats.expiredTickets}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center">
                <XCircle className="h-4 w-4 text-gray-600" />
                <div className="ml-2">
                  <p className="text-sm font-medium">Cancelled</p>
                  <p className="text-2xl font-bold">{stats.cancelledTickets}</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center">
                <Activity className="h-4 w-4 text-purple-600" />
                <div className="ml-2">
                  <p className="text-sm font-medium">Today's Scans</p>
                  <p className="text-2xl font-bold">{stats.todayScans}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      <Tabs defaultValue="validate" className="space-y-4">
        <TabsList>
          <TabsTrigger value="validate">Manual Validation</TabsTrigger>
          <TabsTrigger value="tickets">All Tickets</TabsTrigger>
        </TabsList>

        {/* Manual Validation Tab */}
        <TabsContent value="validate" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <QrCode className="h-5 w-5 mr-2" />
                Manual Ticket Validation
              </CardTitle>
              <CardDescription>
                Enter a QR code manually to validate a ticket
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex space-x-2">
                <div className="flex-1">
                  <Label htmlFor="qr-code">QR Code</Label>
                  <Input
                    id="qr-code"
                    placeholder="Enter QR code (e.g., QR123456789)"
                    value={manualScanCode}
                    onChange={(e) => setManualScanCode(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && handleManualScan()}
                  />
                </div>
                <div className="flex items-end">
                  <Button onClick={handleManualScan}>
                    <Scan className="h-4 w-4 mr-2" />
                    Validate
                  </Button>
                </div>
              </div>

              {/* Scan Result */}
              {scanResult && (
                <Alert variant={scanResult.success ? 'default' : 'destructive'}>
                  {scanResult.success ? (
                    <CheckCircle className="h-4 w-4" />
                  ) : (
                    <AlertTriangle className="h-4 w-4" />
                  )}
                  <AlertDescription>
                    <div className="space-y-2">
                      <p className="font-medium">{scanResult.message}</p>
                      {scanResult.ticket && (
                        <div className="text-sm">
                          <p><strong>Customer:</strong> {scanResult.ticket.customerName}</p>
                          <p><strong>Event:</strong> {scanResult.ticket.eventName}</p>
                          <p><strong>Tier:</strong> {scanResult.ticket.ticketTier}</p>
                          <p><strong>Order:</strong> {scanResult.ticket.orderNumber}</p>
                        </div>
                      )}
                    </div>
                  </AlertDescription>
                </Alert>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* All Tickets Tab */}
        <TabsContent value="tickets" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>All Tickets</CardTitle>
              <CardDescription>
                View and manage all tickets in the system
              </CardDescription>
            </CardHeader>
            <CardContent>
              {/* Filters */}
              <div className="flex space-x-2 mb-4">
                <div className="flex-1">
                  <div className="relative">
                    <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                      placeholder="Search tickets..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="pl-8"
                    />
                  </div>
                </div>
                <Select value={statusFilter} onValueChange={setStatusFilter}>
                  <SelectTrigger className="w-40">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Status</SelectItem>
                    <SelectItem value="valid">Valid</SelectItem>
                    <SelectItem value="used">Used</SelectItem>
                    <SelectItem value="expired">Expired</SelectItem>
                    <SelectItem value="cancelled">Cancelled</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Tickets Table */}
              <div className="rounded-md border">
                <div className="overflow-x-auto">
                  <table className="w-full">
                    <thead>
                      <tr className="border-b bg-muted/50">
                        <th className="h-12 px-4 text-left align-middle font-medium">QR Code</th>
                        <th className="h-12 px-4 text-left align-middle font-medium">Customer</th>
                        <th className="h-12 px-4 text-left align-middle font-medium">Event</th>
                        <th className="h-12 px-4 text-left align-middle font-medium">Tier</th>
                        <th className="h-12 px-4 text-left align-middle font-medium">Status</th>
                        <th className="h-12 px-4 text-left align-middle font-medium">Price</th>
                        <th className="h-12 px-4 text-left align-middle font-medium">Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {filteredTickets.map((ticket) => (
                        <tr key={ticket.id} className="border-b">
                          <td className="p-4 align-middle">
                            <code className="text-sm bg-muted px-2 py-1 rounded">
                              {ticket.qrCode}
                            </code>
                          </td>
                          <td className="p-4 align-middle">
                            <div>
                              <p className="font-medium">{ticket.customerName}</p>
                              <p className="text-sm text-muted-foreground">{ticket.customerEmail}</p>
                            </div>
                          </td>
                          <td className="p-4 align-middle">
                            <div>
                              <p className="font-medium">{ticket.eventName}</p>
                              <p className="text-sm text-muted-foreground">{ticket.venue}</p>
                            </div>
                          </td>
                          <td className="p-4 align-middle">
                            <Badge variant="outline">{ticket.ticketTier}</Badge>
                          </td>
                          <td className="p-4 align-middle">
                            {getStatusBadge(ticket.status)}
                          </td>
                          <td className="p-4 align-middle">
                            ₦{ticket.price.toLocaleString()}
                          </td>
                          <td className="p-4 align-middle">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleViewTicket(ticket.id)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {filteredTickets.length === 0 && (
                <div className="text-center py-8">
                  <p className="text-muted-foreground">No tickets found matching your criteria</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Ticket Details Dialog */}
      <Dialog open={showTicketDialog} onOpenChange={setShowTicketDialog}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Ticket Details</DialogTitle>
            <DialogDescription>
              Complete information about this ticket
            </DialogDescription>
          </DialogHeader>
          {selectedTicket && (
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>QR Code</Label>
                  <code className="block text-sm bg-muted p-2 rounded mt-1">
                    {selectedTicket.qrCode}
                  </code>
                </div>
                <div>
                  <Label>Status</Label>
                  <div className="mt-1">
                    {getStatusBadge(selectedTicket.status)}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Customer Name</Label>
                  <p className="text-sm mt-1">{selectedTicket.customerName}</p>
                </div>
                <div>
                  <Label>Email</Label>
                  <p className="text-sm mt-1">{selectedTicket.customerEmail}</p>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Event</Label>
                  <p className="text-sm mt-1">{selectedTicket.eventName}</p>
                </div>
                <div>
                  <Label>Venue</Label>
                  <p className="text-sm mt-1">{selectedTicket.venue}</p>
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div>
                  <Label>Ticket Tier</Label>
                  <p className="text-sm mt-1">{selectedTicket.ticketTier}</p>
                </div>
                <div>
                  <Label>Price</Label>
                  <p className="text-sm mt-1">₦{selectedTicket.price.toLocaleString()}</p>
                </div>
                <div>
                  <Label>Order Number</Label>
                  <p className="text-sm mt-1">{selectedTicket.orderNumber}</p>
                </div>
              </div>

              {selectedTicket.scanHistory.length > 0 && (
                <div>
                  <Label>Scan History</Label>
                  <div className="mt-2 space-y-2">
                    {selectedTicket.scanHistory.map((scan) => (
                      <div key={scan.id} className="bg-muted p-3 rounded text-sm">
                        <div className="flex justify-between items-start">
                          <div>
                            <p><strong>Scanned by:</strong> {scan.scannedBy}</p>
                            <p><strong>Location:</strong> {scan.scannerLocation}</p>
                            <p><strong>Device:</strong> {scan.scannerDevice}</p>
                          </div>
                          <div className="text-right">
                            <p className="text-muted-foreground">
                              {new Date(scan.scannedAt).toLocaleString()}
                            </p>
                            <Badge variant={scan.result === 'success' ? 'default' : 'destructive'}>
                              {scan.result}
                            </Badge>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <div className="flex justify-end space-x-2">
                {selectedTicket.status === 'valid' && (
                  <Button
                    variant="destructive"
                    onClick={() => handleInvalidateTicket(selectedTicket.id)}
                  >
                    Invalidate Ticket
                  </Button>
                )}
                <Button variant="outline" onClick={() => setShowTicketDialog(false)}>
                  Close
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default AdminTicketValidationPage


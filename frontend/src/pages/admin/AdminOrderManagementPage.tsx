import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { 
  ShoppingCart,
  Search,
  Filter,
  Eye,
  Edit,
  Download,
  RefreshCw,
  Clock,
  User,
  Calendar,
  DollarSign,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Package,
  CreditCard,
  Mail,
  Phone,
  MapPin,
  Ticket,
  MoreHorizontal
} from 'lucide-react'

// TypeScript interfaces
interface Order {
  id: string
  orderCode: string
  customerName: string
  customerEmail: string
  customerPhone?: string
  eventName: string
  ticketTier: string
  ticketQuantity: number
  totalAmount: number
  status: 'pending' | 'confirmed' | 'completed' | 'cancelled' | 'refunded'
  paymentStatus: 'pending' | 'paid' | 'failed' | 'refunded'
  paymentMethod?: string
  paymentReference?: string
  paymentDate?: string
  createdAt: string
}

interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}

interface OrdersQueryParams {
  page?: number
  limit?: number
  status?: string
  paymentStatus?: string
  search?: string
}

const AdminOrderManagementPage: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([])
  const [filteredOrders, setFilteredOrders] = useState<Order[]>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [searchTerm, setSearchTerm] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [paymentFilter, setPaymentFilter] = useState<string>('all')
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const [showOrderDialog, setShowOrderDialog] = useState<boolean>(false)
  const [isExporting, setIsExporting] = useState<boolean>(false)

  useEffect(() => {
    fetchOrders()
  }, [])

  useEffect(() => {
    filterOrders()
  }, [orders, searchTerm, statusFilter, paymentFilter])

  const fetchOrders = async (): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/orders', {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: PaginatedResponse<Order> } = await response.json()
        setOrders(result.data.data || [])
      }
    } catch (error) {
      console.error('Failed to fetch orders:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const filterOrders = (): void => {
    let filtered = orders

    if (searchTerm) {
      filtered = filtered.filter(order => 
        order.orderCode.toLowerCase().includes(searchTerm.toLowerCase()) ||
        order.customerName.toLowerCase().includes(searchTerm.toLowerCase()) ||
        order.customerEmail.toLowerCase().includes(searchTerm.toLowerCase()) ||
        order.eventName.toLowerCase().includes(searchTerm.toLowerCase())
      )
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter(order => order.status === statusFilter)
    }

    if (paymentFilter !== 'all') {
      filtered = filtered.filter(order => order.paymentStatus === paymentFilter)
    }

    setFilteredOrders(filtered)
  }

  const handleViewOrder = async (orderId: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`http://localhost:8080/v1/admin/orders/${orderId}`, {
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        }
      })
      
      if (response.ok) {
        const result: { data: Order } = await response.json()
        setSelectedOrder(result.data)
        setShowOrderDialog(true)
      }
    } catch (error) {
      console.error('Failed to fetch order details:', error)
    }
  }

  const handleUpdateOrderStatus = async (orderId: string, newStatus: string): Promise<void> => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch(`http://localhost:8080/v1/admin/orders/${orderId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${adminToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ status: newStatus })
      })
      
      if (response.ok) {
        fetchOrders() // Refresh the list
        alert('Order status updated successfully!')
      } else {
        alert('Failed to update order status')
      }
    } catch (error) {
      console.error('Failed to update order status:', error)
      alert('An error occurred while updating order status.')
    }
  }

  // NEW: Export Orders functionality
  const handleExportOrders = async (): Promise<void> => {
    setIsExporting(true)
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/orders/export', {
        headers: {
          'Authorization': `Bearer ${adminToken}`
        }
      })
      
      if (response.ok) {
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `orders-export-${new Date().toISOString().split('T')[0]}.csv`
        a.click()
        window.URL.revokeObjectURL(url)
        alert('Orders exported successfully!')
      } else {
        alert('Failed to export orders')
      }
    } catch (error) {
      console.error('Export failed:', error)
      alert('An error occurred while exporting orders.')
    } finally {
      setIsExporting(false)
    }
  }

  // NEW: Refund Order functionality
  const handleRefundOrder = async (orderId: string): Promise<void> => {
    if (confirm('Are you sure you want to refund this order? This action cannot be undone.')) {
      try {
        const adminToken = localStorage.getItem('adminToken')
        const response = await fetch(`http://localhost:8080/v1/admin/orders/${orderId}/refund`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${adminToken}`,
            'Content-Type': 'application/json'
          }
        })
        
        if (response.ok) {
          fetchOrders()
          alert('Order refunded successfully!')
        } else {
          alert('Failed to refund order')
        }
      } catch (error) {
        console.error('Failed to refund order:', error)
        alert('An error occurred while processing the refund.')
      }
    }
  }

  // NEW: Cancel Order functionality
  const handleCancelOrder = async (orderId: string): Promise<void> => {
    if (confirm('Are you sure you want to cancel this order?')) {
      try {
        const adminToken = localStorage.getItem('adminToken')
        const response = await fetch(`http://localhost:8080/v1/admin/orders/${orderId}/cancel`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${adminToken}`,
            'Content-Type': 'application/json'
          }
        })
        
        if (response.ok) {
          fetchOrders()
          alert('Order cancelled successfully!')
        } else {
          alert('Failed to cancel order')
        }
      } catch (error) {
        console.error('Failed to cancel order:', error)
        alert('An error occurred while cancelling the order.')
      }
    }
  }

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: 'NGN'
    }).format(amount)
  }

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString('en-NG', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getStatusBadge = (status: string): JSX.Element => {
    const statusConfig: Record<string, { color: string; label: string; icon: React.ComponentType<any> }> = {
      pending: { color: 'bg-yellow-100 text-yellow-800', label: 'Pending', icon: Clock },
      confirmed: { color: 'bg-blue-100 text-blue-800', label: 'Confirmed', icon: CheckCircle },
      completed: { color: 'bg-green-100 text-green-800', label: 'Completed', icon: CheckCircle },
      cancelled: { color: 'bg-red-100 text-red-800', label: 'Cancelled', icon: XCircle },
      refunded: { color: 'bg-gray-100 text-gray-800', label: 'Refunded', icon: RefreshCw }
    }
    
    const config = statusConfig[status] || { color: 'bg-gray-100 text-gray-800', label: status, icon: AlertTriangle }
    const Icon = config.icon
    
    return (
      <Badge className={`${config.color} flex items-center gap-1`}>
        <Icon className="h-3 w-3" />
        {config.label}
      </Badge>
    )
  }

  const getPaymentStatusBadge = (status: string): JSX.Element => {
    const statusConfig: Record<string, { color: string; label: string }> = {
      pending: { color: 'bg-yellow-100 text-yellow-800', label: 'Pending' },
      paid: { color: 'bg-green-100 text-green-800', label: 'Paid' },
      failed: { color: 'bg-red-100 text-red-800', label: 'Failed' },
      refunded: { color: 'bg-gray-100 text-gray-800', label: 'Refunded' }
    }
    
    const config = statusConfig[status] || { color: 'bg-gray-100 text-gray-800', label: status }
    return (
      <Badge className={config.color}>
        {config.label}
      </Badge>
    )
  }

  const getOrderStats = () => {
    const total = orders.length
    const pending = orders.filter(o => o.status === 'pending').length
    const completed = orders.filter(o => o.status === 'completed').length
    const cancelled = orders.filter(o => o.status === 'cancelled').length
    const totalRevenue = orders
      .filter(o => o.paymentStatus === 'paid')
      .reduce((sum, o) => sum + o.totalAmount, 0)
    
    return { total, pending, completed, cancelled, totalRevenue }
  }

  const stats = getOrderStats()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-2">Loading orders...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Order Management</h1>
          <p className="text-muted-foreground">
            Manage customer orders and process payments
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" onClick={fetchOrders}>
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button 
            variant="outline" 
            onClick={handleExportOrders}
            disabled={isExporting}
          >
            <Download className="h-4 w-4 mr-2" />
            {isExporting ? 'Exporting...' : 'Export Orders'}
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Orders</CardTitle>
            <ShoppingCart className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
            <p className="text-xs text-muted-foreground">
              All time orders
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-600">{stats.pending}</div>
            <p className="text-xs text-muted-foreground">
              Awaiting processing
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Completed</CardTitle>
            <CheckCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{stats.completed}</div>
            <p className="text-xs text-muted-foreground">
              Successfully completed
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cancelled</CardTitle>
            <XCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{stats.cancelled}</div>
            <p className="text-xs text-muted-foreground">
              Cancelled orders
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{formatCurrency(stats.totalRevenue)}</div>
            <p className="text-xs text-muted-foreground">
              From paid orders
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center space-x-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search orders by code, customer, or event..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="confirmed">Confirmed</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="cancelled">Cancelled</SelectItem>
                <SelectItem value="refunded">Refunded</SelectItem>
              </SelectContent>
            </Select>
            <Select value={paymentFilter} onValueChange={setPaymentFilter}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Filter by payment" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Payments</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="paid">Paid</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
                <SelectItem value="refunded">Refunded</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Orders Table */}
      <Card>
        <CardHeader>
          <CardTitle>Orders ({filteredOrders.length})</CardTitle>
          <CardDescription>
            Manage customer orders and track payments
          </CardDescription>
        </CardHeader>
        <CardContent className="p-0">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="border-b">
                <tr>
                  <th className="text-left p-4">Order</th>
                  <th className="text-left p-4">Customer</th>
                  <th className="text-left p-4">Event</th>
                  <th className="text-left p-4">Amount</th>
                  <th className="text-left p-4">Status</th>
                  <th className="text-left p-4">Payment</th>
                  <th className="text-left p-4">Date</th>
                  <th className="text-left p-4">Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredOrders.map((order) => (
                  <tr key={order.id} className="border-b hover:bg-gray-50">
                    <td className="p-4">
                      <div>
                        <div className="font-medium">{order.orderCode}</div>
                        <div className="text-sm text-muted-foreground">
                          {order.ticketQuantity} ticket{order.ticketQuantity > 1 ? 's' : ''}
                        </div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="space-y-1">
                        <div className="font-medium">{order.customerName}</div>
                        <div className="text-sm text-muted-foreground flex items-center">
                          <Mail className="h-3 w-3 mr-1" />
                          {order.customerEmail}
                        </div>
                        {order.customerPhone && (
                          <div className="text-sm text-muted-foreground flex items-center">
                            <Phone className="h-3 w-3 mr-1" />
                            {order.customerPhone}
                          </div>
                        )}
                      </div>
                    </td>
                    <td className="p-4">
                      <div>
                        <div className="font-medium">{order.eventName}</div>
                        <div className="text-sm text-muted-foreground">{order.ticketTier}</div>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="font-medium">{formatCurrency(order.totalAmount)}</div>
                    </td>
                    <td className="p-4">{getStatusBadge(order.status)}</td>
                    <td className="p-4">{getPaymentStatusBadge(order.paymentStatus)}</td>
                    <td className="p-4">
                      <div className="text-sm">
                        {formatDate(order.createdAt)}
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center space-x-2">
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => handleViewOrder(order.id)}
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Select onValueChange={(value) => {
                          if (value === 'refund') {
                            handleRefundOrder(order.id)
                          } else if (value === 'cancel') {
                            handleCancelOrder(order.id)
                          } else {
                            handleUpdateOrderStatus(order.id, value)
                          }
                        }}>
                          <SelectTrigger className="w-8 h-8 p-0">
                            <MoreHorizontal className="h-4 w-4" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="confirmed">Mark Confirmed</SelectItem>
                            <SelectItem value="completed">Mark Completed</SelectItem>
                            <SelectItem value="cancel">Cancel Order</SelectItem>
                            {order.paymentStatus === 'paid' && (
                              <SelectItem value="refund">Process Refund</SelectItem>
                            )}
                          </SelectContent>
                        </Select>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      {/* Order Details Dialog */}
      <Dialog open={showOrderDialog} onOpenChange={setShowOrderDialog}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Order Details</DialogTitle>
            <DialogDescription>
              Complete information about the selected order
            </DialogDescription>
          </DialogHeader>
          
          {selectedOrder && (
            <div className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-lg font-semibold">Order {selectedOrder.orderCode}</h3>
                  <p className="text-muted-foreground">{formatDate(selectedOrder.createdAt)}</p>
                </div>
                <div className="text-right">
                  {getStatusBadge(selectedOrder.status)}
                  <div className="mt-1">
                    {getPaymentStatusBadge(selectedOrder.paymentStatus)}
                  </div>
                </div>
              </div>

              <Tabs defaultValue="details" className="w-full">
                <TabsList className="grid w-full grid-cols-3">
                  <TabsTrigger value="details">Order Details</TabsTrigger>
                  <TabsTrigger value="customer">Customer Info</TabsTrigger>
                  <TabsTrigger value="payment">Payment Info</TabsTrigger>
                </TabsList>
                
                <TabsContent value="details" className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <Label>Event</Label>
                      <p className="text-sm font-medium">{selectedOrder.eventName}</p>
                    </div>
                    <div>
                      <Label>Ticket Tier</Label>
                      <p className="text-sm">{selectedOrder.ticketTier}</p>
                    </div>
                    <div>
                      <Label>Quantity</Label>
                      <p className="text-sm">{selectedOrder.ticketQuantity} ticket(s)</p>
                    </div>
                    <div>
                      <Label>Total Amount</Label>
                      <p className="text-sm font-medium">{formatCurrency(selectedOrder.totalAmount)}</p>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="customer" className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <Label>Name</Label>
                      <p className="text-sm">{selectedOrder.customerName}</p>
                    </div>
                    <div>
                      <Label>Email</Label>
                      <p className="text-sm">{selectedOrder.customerEmail}</p>
                    </div>
                    <div>
                      <Label>Phone</Label>
                      <p className="text-sm">{selectedOrder.customerPhone || 'Not provided'}</p>
                    </div>
                  </div>
                </TabsContent>
                
                <TabsContent value="payment" className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <Label>Payment Method</Label>
                      <p className="text-sm">{selectedOrder.paymentMethod || 'Not specified'}</p>
                    </div>
                    <div>
                      <Label>Payment Reference</Label>
                      <p className="text-sm">{selectedOrder.paymentReference || 'Not available'}</p>
                    </div>
                    <div>
                      <Label>Payment Date</Label>
                      <p className="text-sm">
                        {selectedOrder.paymentDate ? formatDate(selectedOrder.paymentDate) : 'Not paid'}
                      </p>
                    </div>
                  </div>
                </TabsContent>
              </Tabs>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  )
}

export default AdminOrderManagementPage


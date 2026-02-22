import React, { useState, useEffect } from 'react'
import { useAuth } from '../../contexts/AuthContext'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { analyticsAPI } from '@/services/api'
import { DashboardStats, EventStats, SalesAnalytics, UserAnalytics, ApiResponse } from '@/types/api'
import { 
  BarChart3, 
  TrendingUp, 
  Users, 
  Calendar, 
  ShoppingCart, 
  DollarSign,
  Download,
  RefreshCw,
  Activity,
  Target,
  Percent,
  Clock
} from 'lucide-react'

const AdminAnalyticsPage: React.FC = () => {
  const { admin, hasPermission } = useAuth()
  const [dashboardStats, setDashboardStats] = useState<DashboardStats | null>(null)
  const [eventStats, setEventStats] = useState<EventStats[]>([])
  const [salesAnalytics, setSalesAnalytics] = useState<SalesAnalytics | null>(null)
  const [userAnalytics, setUserAnalytics] = useState<UserAnalytics | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState('overview')

  useEffect(() => {
    fetchAnalyticsData()
  }, [])

  const fetchAnalyticsData = async () => {
    try {
      setIsLoading(true)
      setError(null)

      // Fetch dashboard overview
      const dashboardResponse: ApiResponse<DashboardStats> = await analyticsAPI.getDashboard()
      if (dashboardResponse.success && dashboardResponse.data) {
        setDashboardStats(dashboardResponse.data)
      }

      // Fetch event analytics
      try {
        const eventResponse: ApiResponse<EventStats[]> = await analyticsAPI.getEventAnalytics()
        if (eventResponse.success && eventResponse.data) {
          setEventStats(eventResponse.data)
        }
      } catch (eventError) {
        console.warn('Event analytics not available:', eventError)
      }

      // Fetch sales analytics
      try {
        const salesResponse: ApiResponse<SalesAnalytics> = await analyticsAPI.getSalesAnalytics()
        if (salesResponse.success && salesResponse.data) {
          setSalesAnalytics(salesResponse.data)
        }
      } catch (salesError) {
        console.warn('Sales analytics not available:', salesError)
      }

      // Fetch user analytics
      try {
        const userResponse: ApiResponse<UserAnalytics> = await analyticsAPI.getUserAnalytics()
        if (userResponse.success && userResponse.data) {
          setUserAnalytics(userResponse.data)
        }
      } catch (userError) {
        console.warn('User analytics not available:', userError)
      }

    } catch (error) {
      console.error('Error fetching analytics:', error)
      setError(error instanceof Error ? error.message : 'Failed to load analytics data')
    } finally {
      setIsLoading(false)
    }
  }

  const formatNumber = (num: number): string => {
    return new Intl.NumberFormat().format(num)
  }

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: 'NGN'
    }).format(amount)
  }

  const formatPercentage = (value: number): string => {
    return `${value.toFixed(1)}%`
  }

  const handleExportAnalytics = async () => {
    try {
      const adminToken = localStorage.getItem('adminToken')
      const response = await fetch('http://localhost:8080/v1/admin/export/events', {
        headers: { 'Authorization': `Bearer ${adminToken}` }
      })
      
      if (response.ok) {
        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `analytics-export-${new Date().toISOString().split('T')[0]}.csv`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
      } else {
        console.error('Export failed:', response.statusText)
      }
    } catch (error) {
      console.error('Export error:', error)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-2">Loading analytics...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-red-600">Error Loading Analytics</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-600 mb-4">{error}</p>
            <Button onClick={fetchAnalyticsData} className="w-full">
              Try Again
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Analytics</h1>
          <p className="text-gray-600">Platform performance and insights</p>
        </div>
        <div className="flex space-x-2">
          <Button onClick={fetchAnalyticsData} variant="outline">
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          {hasPermission('ANALYTICS_EXPORT') && (
            <Button onClick={handleExportAnalytics}>
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          )}
        </div>
      </div>

      {/* Analytics Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="events">Events</TabsTrigger>
          <TabsTrigger value="sales">Sales</TabsTrigger>
          <TabsTrigger value="users">Users</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          {dashboardStats && (
            <>
              {/* Key Metrics */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
                    <DollarSign className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{formatCurrency(dashboardStats.totalRevenue)}</div>
                    <p className="text-xs text-muted-foreground">
                      {formatCurrency(dashboardStats.revenueThisMonth)} this month
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Total Orders</CardTitle>
                    <ShoppingCart className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{formatNumber(dashboardStats.totalOrders)}</div>
                    <p className="text-xs text-muted-foreground">
                      {formatNumber(dashboardStats.ordersThisMonth)} this month
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Events</CardTitle>
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{formatNumber(dashboardStats.totalEvents)}</div>
                    <p className="text-xs text-muted-foreground">
                      {formatNumber(dashboardStats.activeEvents)} active
                    </p>
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Tickets Sold</CardTitle>
                    <Activity className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{formatNumber(dashboardStats.ticketsSold)}</div>
                    <p className="text-xs text-muted-foreground">
                      {formatNumber(dashboardStats.ticketsScanned)} scanned
                    </p>
                  </CardContent>
                </Card>
              </div>

              {/* Performance Metrics */}
              {dashboardStats.performanceMetrics && (
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center">
                      <Target className="h-5 w-5 mr-2" />
                      Performance Metrics
                    </CardTitle>
                    <CardDescription>
                      Key performance indicators for the platform
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                      <div className="text-center">
                        <div className="text-3xl font-bold text-purple-600">
                          {formatPercentage(dashboardStats.performanceMetrics.eventConversionRate)}
                        </div>
                        <p className="text-sm text-gray-500 mt-1">Event Conversion Rate</p>
                      </div>
                      <div className="text-center">
                        <div className="text-3xl font-bold text-green-600">
                          {formatCurrency(dashboardStats.performanceMetrics.averageOrderValue)}
                        </div>
                        <p className="text-sm text-gray-500 mt-1">Average Order Value</p>
                      </div>
                      <div className="text-center">
                        <div className="text-3xl font-bold text-blue-600">
                          {formatPercentage(dashboardStats.performanceMetrics.customerSatisfaction)}
                        </div>
                        <p className="text-sm text-gray-500 mt-1">Customer Satisfaction</p>
                      </div>
                      <div className="text-center">
                        <div className="text-3xl font-bold text-orange-600">
                          {formatPercentage(dashboardStats.performanceMetrics.platformUptime)}
                        </div>
                        <p className="text-sm text-gray-500 mt-1">Platform Uptime</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              )}

              {/* Top Events */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <TrendingUp className="h-5 w-5 mr-2" />
                    Top Performing Events
                  </CardTitle>
                  <CardDescription>
                    Events ranked by revenue performance
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  {dashboardStats.topEvents && dashboardStats.topEvents.length > 0 ? (
                    <div className="space-y-4">
                      {dashboardStats.topEvents.map((event, index) => (
                        <div key={event.eventId} className="flex items-center justify-between p-4 border rounded-lg">
                          <div className="flex items-center space-x-4">
                            <div className="flex-shrink-0">
                              <div className="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center">
                                <span className="text-lg font-bold text-purple-600">#{index + 1}</span>
                              </div>
                            </div>
                            <div>
                              <h3 className="font-medium text-gray-900">{event.eventName}</h3>
                              <p className="text-sm text-gray-500">{formatNumber(event.ticketsSold)} tickets sold</p>
                            </div>
                          </div>
                          <div className="text-right">
                            <div className="text-lg font-bold text-gray-900">{formatCurrency(event.revenue)}</div>
                            <Badge variant="secondary" className="mt-1">
                              Top {index + 1}
                            </Badge>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-gray-500 text-center py-8">No event data available</p>
                  )}
                </CardContent>
              </Card>
            </>
          )}
        </TabsContent>

        {/* Events Tab */}
        <TabsContent value="events" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Event Analytics</CardTitle>
              <CardDescription>
                Detailed performance metrics for individual events
              </CardDescription>
            </CardHeader>
            <CardContent>
              {eventStats.length > 0 ? (
                <div className="space-y-4">
                  {eventStats.map((event) => (
                    <div key={event.event_id} className="border rounded-lg p-4">
                      <h3 className="font-medium text-gray-900 mb-2">{event.event_name}</h3>
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <p className="text-gray-500">Revenue</p>
                          <p className="font-medium">{formatCurrency(event.total_revenue)}</p>
                        </div>
                        <div>
                          <p className="text-gray-500">Tickets Sold</p>
                          <p className="font-medium">{formatNumber(event.total_tickets_sold)}</p>
                        </div>
                        <div>
                          <p className="text-gray-500">Orders</p>
                          <p className="font-medium">{formatNumber(event.total_orders)}</p>
                        </div>
                        <div>
                          <p className="text-gray-500">Conversion Rate</p>
                          <p className="font-medium">{formatPercentage(event.conversion_rate)}</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-gray-500 text-center py-8">No event analytics available</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Sales Tab */}
        <TabsContent value="sales" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Sales Analytics</CardTitle>
              <CardDescription>
                Revenue and sales performance metrics
              </CardDescription>
            </CardHeader>
            <CardContent>
              {salesAnalytics ? (
                <div className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-green-600">
                        {formatCurrency(salesAnalytics.total_revenue)}
                      </div>
                      <p className="text-sm text-gray-500">Total Revenue</p>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-blue-600">
                        {formatCurrency(salesAnalytics.average_order_value)}
                      </div>
                      <p className="text-sm text-gray-500">Average Order Value</p>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-purple-600">
                        {formatPercentage(salesAnalytics.conversion_rate)}
                      </div>
                      <p className="text-sm text-gray-500">Conversion Rate</p>
                    </div>
                  </div>
                </div>
              ) : (
                <p className="text-gray-500 text-center py-8">No sales analytics available</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Users Tab */}
        <TabsContent value="users" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>User Analytics</CardTitle>
              <CardDescription>
                User engagement and behavior metrics
              </CardDescription>
            </CardHeader>
            <CardContent>
              {userAnalytics ? (
                <div className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-blue-600">
                        {formatNumber(userAnalytics.total_users)}
                      </div>
                      <p className="text-sm text-gray-500">Total Users</p>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-green-600">
                        {formatNumber(userAnalytics.active_users)}
                      </div>
                      <p className="text-sm text-gray-500">Active Users</p>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-purple-600">
                        {formatNumber(userAnalytics.new_users_this_month)}
                      </div>
                      <p className="text-sm text-gray-500">New This Month</p>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-orange-600">
                        {formatPercentage(userAnalytics.user_retention_rate)}
                      </div>
                      <p className="text-sm text-gray-500">Retention Rate</p>
                    </div>
                  </div>
                </div>
              ) : (
                <p className="text-gray-500 text-center py-8">No user analytics available</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

export default AdminAnalyticsPage

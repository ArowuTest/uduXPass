/*
 * AdminDashboard — uduXPass Design System
 * Dark navy/amber, Syne headings, Inter body
 * FIXED: admin?.firstName (not first_name), graceful error state, brand colors
 */
import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { analyticsAPI } from '@/services/api'
import { DashboardStats, ApiResponse } from '@/types/api'
import {
  Users, Calendar, ShoppingCart, TrendingUp, Activity,
  Ticket, Scan, Shield, CreditCard, Building, Plus, Eye, RefreshCw,
  ArrowUpRight, BarChart3
} from 'lucide-react'

const AdminDashboard: React.FC = () => {
  const navigate = useNavigate()
  const { admin, hasPermission } = useAuth()
  const [stats, setStats] = useState<DashboardStats>({
    total_events: 0, active_events: 0, total_orders: 0, total_revenue: 0,
    total_tickets_sold: 0, total_tickets_scanned: 0,
    revenue_this_month: 0, orders_this_month: 0,
    top_events: [], recent_orders: []
  })
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => { fetchDashboardStats() }, [])

  const fetchDashboardStats = async () => {
    try {
      setIsLoading(true)
      setError(null)
      const response: ApiResponse<DashboardStats> = await analyticsAPI.getDashboard()
      if (response.success && response.data) {
        setStats(response.data)
      } else {
        throw new Error(response.error || 'Failed to fetch dashboard stats')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dashboard data')
    } finally {
      setIsLoading(false)
    }
  }

  const fmt = (num: number) => isNaN(num) ? '0' : new Intl.NumberFormat().format(num)
  const fmtCurrency = (amount: number) => {
    if (isNaN(amount)) return '₦0'
    return new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN', minimumFractionDigits: 0 }).format(amount)
  }

  const statCards = [
    { label: 'Total Events', value: fmt(stats.total_events), sub: `${fmt(stats.active_events)} active`, icon: Calendar, color: '#F59E0B' },
    { label: 'Total Orders', value: fmt(stats.total_orders), sub: `${fmt(stats.orders_this_month)} this month`, icon: ShoppingCart, color: '#10b981' },
    { label: 'Total Revenue', value: fmtCurrency(stats.total_revenue), sub: `${fmtCurrency(stats.revenue_this_month)} this month`, icon: TrendingUp, color: '#8b5cf6' },
    { label: 'Tickets Sold', value: fmt(stats.total_tickets_sold), sub: `${fmt(stats.total_tickets_scanned)} scanned`, icon: Ticket, color: '#3b82f6' },
  ]

  const quickActions = [
    { label: 'Create Event', icon: Plus, path: '/admin/events/create', perm: 'events_create' },
    { label: 'View Orders', icon: ShoppingCart, path: '/admin/orders', perm: 'orders_view' },
    { label: 'Manage Users', icon: Users, path: '/admin/users', perm: 'users_view' },
    { label: 'Analytics', icon: BarChart3, path: '/admin/analytics', perm: 'analytics_view' },
    { label: 'Scanner Users', icon: Scan, path: '/admin/scanner-users', perm: 'scanners_view' },
    { label: 'Payments', icon: CreditCard, path: '/admin/payments', perm: 'orders_view' },
  ].filter(a => hasPermission(a.perm))

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="flex flex-col items-center gap-4">
          <div className="w-10 h-10 rounded-full border-2 border-t-transparent animate-spin" style={{ borderColor: 'var(--brand-amber)', borderTopColor: 'transparent' }} />
          <p className="text-sm" style={{ color: '#64748b' }}>Loading dashboard...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center p-8 rounded-2xl max-w-sm w-full"
          style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.08)' }}>
          <div className="w-12 h-12 rounded-xl flex items-center justify-center mx-auto mb-4"
            style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)' }}>
            <Activity className="w-6 h-6" style={{ color: '#f87171' }} />
          </div>
          <h3 className="text-lg font-bold mb-2" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Dashboard Unavailable</h3>
          <p className="text-sm mb-6" style={{ color: '#64748b' }}>Could not load dashboard statistics. The analytics service may be temporarily unavailable.</p>
          <Button onClick={fetchDashboardStats} className="w-full"
            style={{ background: 'var(--brand-amber)', color: '#0f1729', fontWeight: 700, fontFamily: 'var(--font-display)' }}>
            <RefreshCw className="w-4 h-4 mr-2" /> Try Again
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8 animate-fade-in">
      {/* Header */}
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p className="text-xs font-semibold tracking-widest uppercase mb-1" style={{ color: 'var(--brand-amber)', fontFamily: 'var(--font-display)' }}>
            Overview
          </p>
          <h1 className="text-3xl font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
            Dashboard
          </h1>
          <p className="text-sm mt-1" style={{ color: '#64748b' }}>
            Welcome back, <span style={{ color: '#94a3b8' }}>{admin?.firstName || 'Admin'}</span>
          </p>
        </div>
        <Button onClick={fetchDashboardStats} variant="outline" size="sm" className="gap-2"
          style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8', background: 'transparent' }}>
          <RefreshCw className="w-4 h-4" /> Refresh
        </Button>
      </div>

      {/* Stat Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        {statCards.map(card => {
          const Icon = card.icon
          return (
            <div key={card.label} className="p-5 rounded-2xl"
              style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.07)', boxShadow: '0 4px 24px rgba(0,0,0,0.2)' }}>
              <div className="flex items-start justify-between mb-4">
                <div className="w-10 h-10 rounded-xl flex items-center justify-center"
                  style={{ background: `${card.color}18`, border: `1px solid ${card.color}30` }}>
                  <Icon className="w-5 h-5" style={{ color: card.color }} />
                </div>
                <ArrowUpRight className="w-4 h-4" style={{ color: '#334155' }} />
              </div>
              <p className="text-2xl font-bold mb-1" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>{card.value}</p>
              <p className="text-xs font-medium mb-0.5" style={{ color: '#64748b' }}>{card.label}</p>
              <p className="text-xs" style={{ color: '#475569' }}>{card.sub}</p>
            </div>
          )
        })}
      </div>

      {/* Quick Actions */}
      {quickActions.length > 0 && (
        <div>
          <h2 className="text-lg font-bold mb-4" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Quick Actions</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
            {quickActions.map(action => {
              const Icon = action.icon
              return (
                <button key={action.label} onClick={() => navigate(action.path)}
                  className="flex flex-col items-center gap-2 p-4 rounded-xl transition-all duration-200 hover:-translate-y-0.5 cursor-pointer"
                  style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.07)' }}>
                  <div className="w-9 h-9 rounded-lg flex items-center justify-center"
                    style={{ background: 'rgba(245,158,11,0.1)', border: '1px solid rgba(245,158,11,0.2)' }}>
                    <Icon className="w-4 h-4" style={{ color: 'var(--brand-amber)' }} />
                  </div>
                  <span className="text-xs font-medium text-center" style={{ color: '#94a3b8' }}>{action.label}</span>
                </button>
              )
            })}
          </div>
        </div>
      )}

      {/* Top Events + Recent Orders */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Events */}
        <div className="rounded-2xl p-6" style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.07)' }}>
          <div className="flex items-center justify-between mb-5">
            <h2 className="text-base font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Top Events</h2>
            <button onClick={() => navigate('/admin/events')} className="text-xs flex items-center gap-1" style={{ color: 'var(--brand-amber)' }}>
              View all <Eye className="w-3 h-3" />
            </button>
          </div>
          {stats.top_events?.length > 0 ? (
            <div className="space-y-3">
              {stats.top_events.slice(0, 5).map((event, i) => (
                <div key={event.eventId} className="flex items-center gap-3 p-3 rounded-xl"
                  style={{ background: 'rgba(255,255,255,0.03)', border: '1px solid rgba(255,255,255,0.05)' }}>
                  <div className="w-7 h-7 rounded-lg flex items-center justify-center text-xs font-bold flex-shrink-0"
                    style={{ background: i === 0 ? 'rgba(245,158,11,0.2)' : 'rgba(255,255,255,0.05)', color: i === 0 ? 'var(--brand-amber)' : '#64748b', fontFamily: 'var(--font-display)' }}>
                    {i + 1}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate" style={{ color: '#f1f5f9' }}>{event.event_name}</p>
                    <p className="text-xs" style={{ color: '#64748b' }}>{fmt(event.tickets_sold)} tickets · {fmtCurrency(event.revenue)}</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <Calendar className="w-8 h-8 mx-auto mb-2" style={{ color: '#334155' }} />
              <p className="text-sm" style={{ color: '#64748b' }}>No event data yet</p>
            </div>
          )}
        </div>

        {/* Recent Orders */}
        <div className="rounded-2xl p-6" style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.07)' }}>
          <div className="flex items-center justify-between mb-5">
            <h2 className="text-base font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Recent Orders</h2>
            <button onClick={() => navigate('/admin/orders')} className="text-xs flex items-center gap-1" style={{ color: 'var(--brand-amber)' }}>
              View all <Eye className="w-3 h-3" />
            </button>
          </div>
          {stats.recent_orders?.length > 0 ? (
            <div className="space-y-3">
              {stats.recent_orders.slice(0, 5).map(order => (
                <div key={order.orderId} className="flex items-center gap-3 p-3 rounded-xl"
                  style={{ background: 'rgba(255,255,255,0.03)', border: '1px solid rgba(255,255,255,0.05)' }}>
                  <div className="w-7 h-7 rounded-lg flex items-center justify-center flex-shrink-0"
                    style={{ background: 'rgba(16,185,129,0.1)' }}>
                    <ShoppingCart className="w-3.5 h-3.5" style={{ color: '#10b981' }} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium truncate" style={{ color: '#f1f5f9' }}>{order.user_name}</p>
                    <p className="text-xs truncate" style={{ color: '#64748b' }}>{order.event_name}</p>
                  </div>
                  <div className="text-right flex-shrink-0">
                    <p className="text-sm font-bold" style={{ color: 'var(--brand-amber)', fontFamily: 'var(--font-display)' }}>{fmtCurrency(order.total_amount)}</p>
                    <Badge className="text-xs px-1.5 py-0.5 mt-0.5"
                      style={{
                        background: order.status === 'completed' ? 'rgba(16,185,129,0.15)' : 'rgba(245,158,11,0.15)',
                        color: order.status === 'completed' ? '#10b981' : 'var(--brand-amber)',
                        border: 'none',
                      }}>
                      {order.status}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8">
              <ShoppingCart className="w-8 h-8 mx-auto mb-2" style={{ color: '#334155' }} />
              <p className="text-sm" style={{ color: '#64748b' }}>No orders yet</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default AdminDashboard

/*
 * AdminLayout — uduXPass Admin
 * Design: Dark navy (#0f1729) sidebar, amber (#F59E0B) accents
 * FIXED: Permission keys now match AuthContext grants (plural form)
 */

import { useState, ReactNode } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import {
  Menu, X, Home, Calendar, Users, ShoppingCart, BarChart3, Settings,
  LogOut, Shield, UserCheck, Ticket, Building, CreditCard, Bell, Globe, Scan, LucideIcon
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'

interface NavigationItem {
  title: string
  path: string
  icon: LucideIcon
  permission: string | null
  roles?: string[]
}
interface AdminLayoutProps { children: ReactNode }

const AdminLayout: React.FC<AdminLayoutProps> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(true)
  const { admin, logout, hasPermission, canAccess } = useAuth()
  const location = useLocation()
  const navigate = useNavigate()

  const handleLogout = (): void => {
    logout()
    navigate('/admin/login')
  }

  const isActivePath = (path: string): boolean =>
    location.pathname === path || location.pathname.startsWith(path + '/')

  // FIXED: Permission keys now match what AuthContext grants to super_admin
  const navigationItems: NavigationItem[] = [
    { title: 'Dashboard',     path: '/admin/dashboard',     icon: Home,         permission: null },
    { title: 'Events',        path: '/admin/events',        icon: Calendar,     permission: 'events_view' },
    { title: 'Orders',        path: '/admin/orders',        icon: ShoppingCart, permission: 'orders_view' },
    { title: 'Users',         path: '/admin/users',         icon: Users,        permission: 'users_view' },
    { title: 'Organizers',    path: '/admin/organizers',    icon: Building,     permission: null },
    { title: 'Tickets',       path: '/admin/tickets',       icon: Ticket,       permission: 'tickets_view' },
    { title: 'Payments',      path: '/admin/payments',      icon: CreditCard,   permission: null },
    { title: 'Analytics',     path: '/admin/analytics',     icon: BarChart3,    permission: 'analytics_view' },
    { title: 'Admin Users',   path: '/admin/admin-users',   icon: UserCheck,    permission: 'admin_create',    roles: ['super_admin'] },
    { title: 'Scanner Users', path: '/admin/scanner-users', icon: Scan,         permission: 'scanners_view',   roles: ['super_admin', 'event_manager'] },
    { title: 'Settings',      path: '/admin/settings',      icon: Settings,     permission: 'settings_update', roles: ['super_admin'] },
  ]

  const filteredNavigation = navigationItems.filter(item => {
    if (!item.permission) return true
    if (item.roles && item.roles.length > 0) return canAccess([item.permission], item.roles)
    return hasPermission(item.permission)
  })

  const initials = `${admin?.firstName?.[0] || ''}${admin?.lastName?.[0] || ''}`.toUpperCase() || 'SA'
  const permCount = admin?.permissions?.includes('*' as never) ? 'All' : `${admin?.permissions?.length || 0}`

  const sidebarBg = '#0f1729'
  const sidebarBorder = 'rgba(255,255,255,0.08)'
  const amber = '#F59E0B'
  const textPrimary = '#f1f5f9'
  const textSecondary = '#94a3b8'

  return (
    <div className="min-h-screen flex" style={{ background: 'var(--bg-primary)' }}>
      {/* Sidebar */}
      <div
        className="flex flex-col transition-all duration-300 flex-shrink-0"
        style={{
          width: isSidebarOpen ? '240px' : '64px',
          background: sidebarBg,
          borderRight: `1px solid ${sidebarBorder}`,
        }}
      >
        {/* Logo Header */}
        <div className="flex items-center justify-between px-4 py-4" style={{ borderBottom: `1px solid ${sidebarBorder}`, minHeight: '64px' }}>
          {isSidebarOpen && (
            <div className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: amber }}>
                <Shield size={16} style={{ color: '#0f1729' }} />
              </div>
              <span className="text-sm font-bold tracking-wide" style={{ color: textPrimary }}>Admin Portal</span>
            </div>
          )}
          <button
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
            className="p-1.5 rounded-lg transition-colors hover:opacity-80"
            style={{ color: textSecondary, marginLeft: isSidebarOpen ? 'auto' : '0' }}
          >
            {isSidebarOpen ? <X size={16} /> : <Menu size={16} />}
          </button>
        </div>

        {/* Admin Profile */}
        {isSidebarOpen && (
          <div className="px-4 py-4" style={{ borderBottom: `1px solid ${sidebarBorder}` }}>
            <div className="flex items-center gap-3">
              <div className="w-9 h-9 rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0"
                style={{ background: amber, color: '#0f1729' }}>
                {initials}
              </div>
              <div className="min-w-0">
                <p className="text-sm font-semibold truncate" style={{ color: textPrimary }}>
                  {admin?.firstName} {admin?.lastName}
                </p>
                <p className="text-xs truncate" style={{ color: textSecondary }}>{admin?.email}</p>
              </div>
            </div>
            <div className="mt-2.5">
              <span className="text-xs font-semibold px-2 py-0.5 rounded-full"
                style={{ background: 'rgba(245,158,11,0.15)', color: amber }}>
                {admin?.role?.replace(/_/g, ' ').toUpperCase() || 'ADMIN'}
              </span>
            </div>
          </div>
        )}

        {/* Navigation */}
        <nav className="flex-1 px-3 py-4 overflow-y-auto">
          <ul className="space-y-1">
            {filteredNavigation.map((item) => {
              const Icon = item.icon
              const isActive = isActivePath(item.path)
              return (
                <li key={item.path}>
                  <Link
                    to={item.path}
                    title={!isSidebarOpen ? item.title : undefined}
                    className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all"
                    style={{
                      background: isActive ? 'rgba(245,158,11,0.15)' : 'transparent',
                      color: isActive ? amber : textSecondary,
                      borderLeft: isActive ? `3px solid ${amber}` : '3px solid transparent',
                    }}
                  >
                    <Icon size={17} className="flex-shrink-0" />
                    {isSidebarOpen && <span>{item.title}</span>}
                  </Link>
                </li>
              )
            })}
          </ul>
        </nav>

        {/* Footer */}
        <div className="px-3 py-4" style={{ borderTop: `1px solid ${sidebarBorder}` }}>
          <Link
            to="/"
            title={!isSidebarOpen ? 'View Site' : undefined}
            className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all hover:opacity-80"
            style={{ color: textSecondary }}
          >
            <Globe size={17} className="flex-shrink-0" />
            {isSidebarOpen && <span>View Site</span>}
          </Link>
          <button
            onClick={handleLogout}
            title={!isSidebarOpen ? 'Logout' : undefined}
            className="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all hover:opacity-80 mt-1"
            style={{ color: '#f87171' }}
          >
            <LogOut size={17} className="flex-shrink-0" />
            {isSidebarOpen && <span>Logout</span>}
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Bar */}
        <header className="flex items-center justify-between px-6 py-4 flex-shrink-0"
          style={{ background: 'var(--bg-elevated)', borderBottom: '1px solid var(--border-color)', minHeight: '64px' }}>
          <h2 className="text-base font-semibold" style={{ color: 'var(--text-primary)' }}>
            uduXPass Administration
          </h2>
          <div className="flex items-center gap-3">
            <button className="p-2 rounded-lg transition-colors hover:opacity-80" style={{ color: 'var(--text-secondary)' }}>
              <Bell size={18} />
            </button>
            <span className="text-xs px-2.5 py-1 rounded-full font-medium"
              style={{ background: 'rgba(245,158,11,0.15)', color: amber }}>
              {permCount} Permissions
            </span>
          </div>
        </header>
        <main className="flex-1 overflow-y-auto p-6" style={{ background: 'var(--brand-navy)' }}>
          {children}
        </main>
      </div>
    </div>
  )
}

export default AdminLayout

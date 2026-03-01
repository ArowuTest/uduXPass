/*
 * AdminLayout — uduXPass Design System
 * Dark navy sidebar, amber accents, Syne font
 * FIXED: Permission keys now match AuthContext (plural: events_view, orders_view, etc.)
 */
import { useState, ReactNode } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import {
  Menu, Home, Calendar, Users, ShoppingCart, BarChart3, Settings, LogOut,
  Shield, UserCheck, Ticket, Building, CreditCard, Scan, ChevronLeft, LucideIcon
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

  const handleLogout = () => { logout(); navigate('/admin/login') }
  const isActive = (path: string) => location.pathname === path || location.pathname.startsWith(path + '/')

  // FIXED: Permission keys match AuthContext (plural underscore format)
  const navigationItems: NavigationItem[] = [
    { title: 'Dashboard',     path: '/admin/dashboard',     icon: Home,         permission: null },
    { title: 'Events',        path: '/admin/events',        icon: Calendar,     permission: 'events_view' },
    { title: 'Orders',        path: '/admin/orders',        icon: ShoppingCart, permission: 'orders_view' },
    { title: 'Users',         path: '/admin/users',         icon: Users,        permission: 'users_view' },
    { title: 'Organizers',    path: '/admin/organizers',    icon: Building,     permission: 'users_view' },
    { title: 'Tickets',       path: '/admin/tickets',       icon: Ticket,       permission: 'tickets_view' },
    { title: 'Payments',      path: '/admin/payments',      icon: CreditCard,   permission: 'orders_view' },
    { title: 'Analytics',     path: '/admin/analytics',     icon: BarChart3,    permission: 'analytics_view' },
    { title: 'Admin Users',   path: '/admin/admin-users',   icon: UserCheck,    permission: 'admin_create',    roles: ['super_admin'] },
    { title: 'Scanner Users', path: '/admin/scanner-users', icon: Scan,         permission: 'scanners_view',   roles: ['super_admin', 'event_manager'] },
    { title: 'Settings',      path: '/admin/settings',      icon: Settings,     permission: 'settings_update', roles: ['super_admin'] },
  ]

  const filteredNavigation = navigationItems.filter(item => {
    if (!item.permission) return true
    if (item.roles?.length) return canAccess([item.permission], item.roles)
    return hasPermission(item.permission)
  })

  const roleLabel = (role: string) => role.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())

  return (
    <div className="flex h-screen overflow-hidden" style={{ background: 'var(--brand-navy)' }}>
      {/* Sidebar */}
      <aside
        className="flex flex-col transition-all duration-300 flex-shrink-0"
        style={{ width: isSidebarOpen ? '240px' : '64px', background: 'var(--brand-surface)', borderRight: '1px solid rgba(255,255,255,0.07)' }}
      >
        {/* Logo */}
        <div className="flex items-center justify-between px-4 py-4" style={{ borderBottom: '1px solid rgba(255,255,255,0.07)', minHeight: '64px' }}>
          {isSidebarOpen ? (
            <>
              <Link to="/admin/dashboard" className="flex items-center gap-2" style={{ textDecoration: 'none' }}>
                <div className="w-7 h-7 rounded-lg flex items-center justify-center" style={{ background: 'var(--brand-amber)' }}>
                  <Shield className="w-4 h-4" style={{ color: '#0f1729' }} />
                </div>
                <span className="text-base font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
                  uduX<span style={{ color: 'var(--brand-amber)' }}>Admin</span>
                </span>
              </Link>
              <button onClick={() => setIsSidebarOpen(false)} className="p-1 rounded" style={{ color: '#64748b' }}>
                <ChevronLeft className="w-4 h-4" />
              </button>
            </>
          ) : (
            <button onClick={() => setIsSidebarOpen(true)} className="w-7 h-7 rounded-lg flex items-center justify-center mx-auto" style={{ background: 'var(--brand-amber)' }}>
              <Shield className="w-4 h-4" style={{ color: '#0f1729' }} />
            </button>
          )}
        </div>

        {/* Nav */}
        <nav className="flex-1 overflow-y-auto py-4 px-2">
          {filteredNavigation.map(item => {
            const Icon = item.icon
            const active = isActive(item.path)
            return (
              <Link key={item.path} to={item.path} style={{ textDecoration: 'none' }}>
                <div
                  className="flex items-center gap-3 px-3 py-2.5 mb-0.5 rounded-lg transition-all duration-150 cursor-pointer"
                  style={{ color: active ? 'var(--brand-amber)' : '#94a3b8', background: active ? 'rgba(245,158,11,0.12)' : 'transparent' }}
                  title={!isSidebarOpen ? item.title : undefined}
                >
                  <Icon className="w-4 h-4 flex-shrink-0" />
                  {isSidebarOpen && <span className="text-sm font-medium">{item.title}</span>}
                  {isSidebarOpen && active && <div className="ml-auto w-1.5 h-1.5 rounded-full" style={{ background: 'var(--brand-amber)' }} />}
                </div>
              </Link>
            )
          })}
        </nav>

        {/* Footer */}
        <div className="p-3" style={{ borderTop: '1px solid rgba(255,255,255,0.07)' }}>
          {isSidebarOpen && (
            <div className="flex items-center gap-3 mb-3 px-2">
              <div className="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0"
                style={{ background: 'var(--brand-amber)', color: '#0f1729' }}>
                {admin?.firstName?.[0]?.toUpperCase() || 'A'}
              </div>
              <div className="min-w-0">
                <p className="text-sm font-semibold truncate" style={{ color: '#f1f5f9', fontFamily: 'var(--font-display)' }}>
                  {admin?.firstName} {admin?.lastName}
                </p>
                <p className="text-xs truncate" style={{ color: '#64748b' }}>{admin?.role ? roleLabel(admin.role) : 'Admin'}</p>
              </div>
            </div>
          )}
          <button onClick={handleLogout} className="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm" style={{ color: '#64748b' }} title="Sign Out">
            <LogOut className="w-4 h-4 flex-shrink-0" />
            {isSidebarOpen && <span>Sign Out</span>}
          </button>
        </div>
      </aside>

      {/* Main */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="flex items-center gap-4 px-6 h-16 flex-shrink-0"
          style={{ background: 'var(--brand-surface)', borderBottom: '1px solid rgba(255,255,255,0.07)' }}>
          {!isSidebarOpen && (
            <button onClick={() => setIsSidebarOpen(true)} className="p-1.5 rounded-lg" style={{ color: '#94a3b8' }}>
              <Menu className="w-5 h-5" />
            </button>
          )}
          <div className="flex-1" />
          <div className="flex items-center gap-3">
            <Badge className="text-xs px-2 py-1 font-semibold"
              style={{ background: 'rgba(245,158,11,0.15)', color: 'var(--brand-amber)', border: '1px solid rgba(245,158,11,0.3)', fontFamily: 'var(--font-display)' }}>
              <Shield className="w-3 h-3 mr-1" />
              {admin?.role ? roleLabel(admin.role) : 'Admin'}
            </Badge>
            <div className="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold"
              style={{ background: 'var(--brand-amber)', color: '#0f1729' }}>
              {admin?.firstName?.[0]?.toUpperCase() || 'A'}
            </div>
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

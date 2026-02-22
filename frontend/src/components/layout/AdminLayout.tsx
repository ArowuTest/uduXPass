import { useState, ReactNode } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { 
  Menu, 
  X, 
  Home, 
  Calendar, 
  Users, 
  ShoppingCart, 
  BarChart3, 
  Settings, 
  LogOut,
  Shield,
  UserCheck,
  Ticket,
  Building,
  CreditCard,
  Bell,
  HelpCircle,
  Scan,
  LucideIcon
} from 'lucide-react'

interface NavigationItem {
  title: string
  path: string
  icon: LucideIcon
  permission: string | null
  roles?: string[]
}

interface AdminLayoutProps {
  children: ReactNode
}

const AdminLayout: React.FC<AdminLayoutProps> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState<boolean>(true)
  const { admin, logout, hasPermission, canAccess } = useAuth()
  const location = useLocation()
  const navigate = useNavigate()

  const handleLogout = (): void => {
    logout()
    navigate('/admin/login')
  }

  const isActivePath = (path: string): boolean => {
    return location.pathname === path || location.pathname.startsWith(path + '/')
  }

  const getRoleColor = (role: string): string => {
    switch (role) {
      case 'super_admin': return 'bg-red-100 text-red-800'
      case 'event_manager': return 'bg-blue-100 text-blue-800'
      case 'support_agent': return 'bg-green-100 text-green-800'
      case 'analyst': return 'bg-purple-100 text-purple-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const navigationItems: NavigationItem[] = [
    {
      title: 'Dashboard',
      path: '/admin/dashboard',
      icon: Home,
      permission: null // Always accessible to admins
    },
    {
      title: 'Events',
      path: '/admin/events',
      icon: Calendar,
      permission: 'event_view'
    },
    {
      title: 'Orders',
      path: '/admin/orders',
      icon: ShoppingCart,
      permission: 'order_view'
    },
    {
      title: 'Users',
      path: '/admin/users',
      icon: Users,
      permission: 'user_view'
    },
    {
      title: 'Organizers',
      path: '/admin/organizers',
      icon: Building,
      permission: 'organizer_view'
    },
    {
      title: 'Tickets',
      path: '/admin/tickets',
      icon: Ticket,
      permission: 'ticket_view'
    },
    {
      title: 'Payments',
      path: '/admin/payments',
      icon: CreditCard,
      permission: 'payment_view'
    },
    {
      title: 'Analytics',
      path: '/admin/analytics',
      icon: BarChart3,
      permission: 'analytics_view'
    },
    {
      title: 'Admin Users',
      path: '/admin/admin-users',
      icon: UserCheck,
      permission: 'admin_create',
      roles: ['super_admin']
    },
    {
      title: 'Scanner Users',
      path: '/admin/scanner-users',
      icon: Scan,
      permission: 'scanner_manage',
      roles: ['super_admin', 'event_manager']
    },
    {
      title: 'Settings',
      path: '/admin/settings',
      icon: Settings,
      permission: 'system_settings',
      roles: ['super_admin']
    }
  ]

  const filteredNavigation = navigationItems.filter(item => {
    if (!item.permission) return true // Always show items without permission requirements
    
    if (item.roles && item.roles.length > 0) {
      return canAccess([item.permission], item.roles)
    }
    
    return hasPermission(item.permission)
  })

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Sidebar */}
      <div className={`${isSidebarOpen ? 'w-64' : 'w-16'} bg-white border-r border-gray-200 transition-all duration-300 flex flex-col`}>
        {/* Header */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            {isSidebarOpen && (
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-gradient-to-br from-purple-600 to-blue-600 rounded-lg flex items-center justify-center">
                  <Shield className="w-5 h-5 text-white" />
                </div>
                <span className="text-lg font-bold text-gray-900">Admin Portal</span>
              </div>
            )}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              className="p-1"
            >
              {isSidebarOpen ? <X className="w-4 h-4" /> : <Menu className="w-4 h-4" />}
            </Button>
          </div>
        </div>

        {/* Admin Info */}
        {isSidebarOpen && (
          <div className="p-4 border-b border-gray-200">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-br from-purple-600 to-blue-600 rounded-full flex items-center justify-center">
                <span className="text-white font-medium">
                  {admin?.firstName?.[0]}{admin?.lastName?.[0]}
                </span>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {admin?.firstName} {admin?.lastName}
                </p>
                <p className="text-xs text-gray-500 truncate">{admin?.email}</p>
              </div>
            </div>
            <div className="mt-2">
              <Badge className={`text-xs ${getRoleColor(admin?.role || '')}`}>
                {admin?.role?.replace('_', ' ').toUpperCase()}
              </Badge>
            </div>
          </div>
        )}

        {/* Navigation */}
        <nav className="flex-1 p-4">
          <ul className="space-y-2">
            {filteredNavigation.map((item) => {
              const Icon = item.icon
              const isActive = isActivePath(item.path)
              
              return (
                <li key={item.path}>
                  <Link
                    to={item.path}
                    className={`flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                      isActive
                        ? 'bg-purple-100 text-purple-700'
                        : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                    }`}
                  >
                    <Icon className="w-5 h-5 flex-shrink-0" />
                    {isSidebarOpen && <span>{item.title}</span>}
                  </Link>
                </li>
              )
            })}
          </ul>
        </nav>

        {/* Footer */}
        <div className="p-4 border-t border-gray-200">
          <div className="space-y-2">
            <Link
              to="/"
              className="flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
            >
              <HelpCircle className="w-5 h-5 flex-shrink-0" />
              {isSidebarOpen && <span>View Site</span>}
            </Link>
            <Button
              variant="ghost"
              className="w-full justify-start px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900"
              onClick={handleLogout}
            >
              <LogOut className="w-5 h-5 flex-shrink-0 mr-3" />
              {isSidebarOpen && <span>Logout</span>}
            </Button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Bar */}
        <header className="bg-white border-b border-gray-200 px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <h2 className="text-lg font-semibold text-gray-900">
                uduXPass Administration
              </h2>
            </div>
            <div className="flex items-center space-x-4">
              <Button variant="ghost" size="sm">
                <Bell className="w-4 h-4" />
              </Button>
              <Badge variant="outline" className="text-xs">
                {admin?.permissions?.includes('*' as any) ? 'All Permissions' : `${admin?.permissions?.length || 0} Permissions`}
              </Badge>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  )
}

export default AdminLayout


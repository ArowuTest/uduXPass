/*
 * Navbar — uduXPass Design System
 * Dark navy (#0f1729) background, amber (#F59E0B) accents
 * No Admin button in public nav — admin access is via /admin/login only
 * Font: Playfair Display (brand) + Inter (body)
 */
import { useState } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { useCart } from '../../contexts/CartContext'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Menu,
  X,
  ShoppingCart,
  User,
  LogOut,
  Ticket,
  Search,
  Calendar,
  Shield,
} from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const Navbar: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const { user, admin, isAuthenticated, isAdmin, logout } = useAuth()
  const { getTotalItems, toggleCart } = useCart()
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  const isActive = (path: string) => location.pathname === path

  // Admin pages use AdminLayout's own header — hide Navbar
  if (location.pathname.startsWith('/admin')) return null

  return (
    <nav
      style={{ background: 'rgba(15,23,41,0.97)', borderBottom: '1px solid rgba(245,158,11,0.15)' }}
      className="sticky top-0 z-50 backdrop-blur-md"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">

          {/* Logo */}
          <Link
            to={isAdmin ? '/admin/dashboard' : '/home'}
            className="flex items-center space-x-2 group"
          >
            <div
              className="w-8 h-8 rounded-lg flex items-center justify-center"
              style={{ background: 'linear-gradient(135deg,#F59E0B,#D97706)' }}
            >
              <Ticket className="w-5 h-5 text-white" />
            </div>
            <span
              className="text-xl font-bold"
              style={{ fontFamily: "'Playfair Display', serif", color: '#F59E0B' }}
            >
              uduXPass
            </span>
          </Link>

          {/* Desktop nav links */}
          {!isAdmin && (
            <div className="hidden md:flex items-center space-x-1">
              {[
                { path: '/home', label: 'Home' },
                { path: '/events', label: 'Events', icon: Calendar },
              ].map(({ path, label, icon: Icon }) => (
                <Link
                  key={path}
                  to={path}
                  className="flex items-center space-x-1 px-4 py-2 rounded-md text-sm font-medium transition-all duration-200"
                  style={{
                    color: isActive(path) ? '#F59E0B' : 'rgba(255,255,255,0.75)',
                    background: isActive(path) ? 'rgba(245,158,11,0.1)' : 'transparent',
                  }}
                  onMouseEnter={e => {
                    if (!isActive(path)) {
                      (e.currentTarget as HTMLElement).style.color = '#F59E0B'
                      ;(e.currentTarget as HTMLElement).style.background = 'rgba(245,158,11,0.08)'
                    }
                  }}
                  onMouseLeave={e => {
                    if (!isActive(path)) {
                      (e.currentTarget as HTMLElement).style.color = 'rgba(255,255,255,0.75)'
                      ;(e.currentTarget as HTMLElement).style.background = 'transparent'
                    }
                  }}
                >
                  {Icon && <Icon className="w-4 h-4" />}
                  <span>{label}</span>
                </Link>
              ))}
            </div>
          )}

          {/* Right side */}
          <div className="flex items-center space-x-2">

            {/* Search */}
            {!isAdmin && (
              <Button variant="ghost" size="sm" className="hidden sm:flex text-white/70 hover:text-amber-400">
                <Search className="w-4 h-4" />
              </Button>
            )}

            {/* Cart */}
            {!isAdmin && (
              <Button
                variant="ghost"
                size="sm"
                onClick={toggleCart}
                className="relative text-white/70 hover:text-amber-400"
              >
                <ShoppingCart className="w-4 h-4" />
                {getTotalItems() > 0 && (
                  <Badge
                    className="absolute -top-2 -right-2 w-5 h-5 flex items-center justify-center p-0 text-xs"
                    style={{ background: '#F59E0B', color: '#0f1729' }}
                  >
                    {getTotalItems()}
                  </Badge>
                )}
              </Button>
            )}

            {/* Admin badge (when logged in as admin) */}
            {isAdmin && (
              <Badge
                className="text-xs"
                style={{ background: 'rgba(245,158,11,0.15)', color: '#F59E0B', border: '1px solid rgba(245,158,11,0.3)' }}
              >
                <Shield className="w-3 h-3 mr-1" />
                Admin
              </Badge>
            )}

            {/* User menu */}
            {isAuthenticated ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="flex items-center space-x-2 text-white/80 hover:text-amber-400"
                  >
                    <div
                      className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold"
                      style={{ background: 'linear-gradient(135deg,#F59E0B,#D97706)', color: '#0f1729' }}
                    >
                      {(isAdmin ? admin?.firstName : user?.firstName)?.[0]?.toUpperCase() || 'U'}
                    </div>
                    <span className="hidden sm:inline text-sm">
                      {isAdmin ? admin?.firstName : user?.firstName}
                    </span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent
                  align="end"
                  className="w-56"
                  style={{ background: '#0f1729', border: '1px solid rgba(245,158,11,0.2)' }}
                >
                  <DropdownMenuLabel>
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium text-white">
                        {isAdmin
                          ? `${admin?.firstName} ${admin?.lastName}`
                          : `${user?.firstName} ${user?.lastName}`}
                      </p>
                      <p className="text-xs text-white/50">
                        {isAdmin ? admin?.email : user?.email}
                      </p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator style={{ background: 'rgba(245,158,11,0.15)' }} />
                  {isAdmin ? (
                    <>
                      <DropdownMenuItem
                        onClick={() => navigate('/admin/dashboard')}
                        className="text-white/80 hover:text-amber-400 cursor-pointer"
                      >
                        <Shield className="mr-2 h-4 w-4" />
                        Admin Dashboard
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => navigate('/home')}
                        className="text-white/80 hover:text-amber-400 cursor-pointer"
                      >
                        <Ticket className="mr-2 h-4 w-4" />
                        View Public Site
                      </DropdownMenuItem>
                    </>
                  ) : (
                    <>
                      <DropdownMenuItem
                        onClick={() => navigate('/profile')}
                        className="text-white/80 hover:text-amber-400 cursor-pointer"
                      >
                        <User className="mr-2 h-4 w-4" />
                        Profile
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => navigate('/profile?tab=orders')}
                        className="text-white/80 hover:text-amber-400 cursor-pointer"
                      >
                        <Ticket className="mr-2 h-4 w-4" />
                        My Orders
                      </DropdownMenuItem>
                    </>
                  )}
                  <DropdownMenuSeparator style={{ background: 'rgba(245,158,11,0.15)' }} />
                  <DropdownMenuItem
                    onClick={handleLogout}
                    className="text-red-400 hover:text-red-300 cursor-pointer"
                  >
                    <LogOut className="mr-2 h-4 w-4" />
                    Log out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <div className="flex items-center space-x-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => navigate('/login')}
                  className="text-white/75 hover:text-amber-400"
                >
                  Sign In
                </Button>
                <Button
                  size="sm"
                  onClick={() => navigate('/register')}
                  style={{ background: 'linear-gradient(135deg,#F59E0B,#D97706)', color: '#0f1729' }}
                  className="font-semibold hover:opacity-90 transition-opacity"
                >
                  Get Tickets
                </Button>
              </div>
            )}

            {/* Mobile menu toggle */}
            {!isAdmin && (
              <Button
                variant="ghost"
                size="sm"
                className="md:hidden text-white/70 hover:text-amber-400"
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              >
                {isMobileMenuOpen ? <X className="w-4 h-4" /> : <Menu className="w-4 h-4" />}
              </Button>
            )}
          </div>
        </div>

        {/* Mobile menu */}
        {!isAdmin && isMobileMenuOpen && (
          <div
            className="md:hidden py-4"
            style={{ borderTop: '1px solid rgba(245,158,11,0.15)' }}
          >
            <div className="flex flex-col space-y-1">
              {[
                { path: '/home', label: 'Home' },
                { path: '/events', label: 'Events' },
              ].map(({ path, label }) => (
                <Link
                  key={path}
                  to={path}
                  className="px-3 py-2 rounded-md text-sm font-medium"
                  style={{ color: isActive(path) ? '#F59E0B' : 'rgba(255,255,255,0.75)' }}
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  {label}
                </Link>
              ))}

              {!isAuthenticated && (
                <div
                  className="flex flex-col space-y-2 pt-4 mt-2"
                  style={{ borderTop: '1px solid rgba(245,158,11,0.15)' }}
                >
                  <Button
                    variant="ghost"
                    className="justify-start text-white/75 hover:text-amber-400"
                    onClick={() => { navigate('/login'); setIsMobileMenuOpen(false) }}
                  >
                    Sign In
                  </Button>
                  <Button
                    className="justify-start font-semibold"
                    style={{ background: 'linear-gradient(135deg,#F59E0B,#D97706)', color: '#0f1729' }}
                    onClick={() => { navigate('/register'); setIsMobileMenuOpen(false) }}
                  >
                    Get Tickets
                  </Button>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  )
}

export default Navbar

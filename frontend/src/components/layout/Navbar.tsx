/*
 * Navbar — uduXPass Design System
 * Dark navy, amber accent, Syne logo — NO admin link in public nav
 */
import { useState } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../../contexts/AuthContext'
import { useCart } from '../../contexts/CartContext'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Menu, X, ShoppingCart, User, LogOut, Ticket, Calendar, ChevronDown } from 'lucide-react'
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem,
  DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const Navbar: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const { user, isAuthenticated, logout } = useAuth()
  const { getTotalItems, toggleCart } = useCart()
  const navigate = useNavigate()
  const location = useLocation()

  if (location.pathname.startsWith('/admin')) return null

  const handleLogout = () => { logout(); navigate('/') }
  const isActive = (path: string) => location.pathname === path || location.pathname.startsWith(path + '/')
  const cartCount = getTotalItems()

  return (
    <nav className="sticky top-0 z-50 border-b" style={{ background: 'rgba(15,23,41,0.96)', backdropFilter: 'blur(12px)', borderColor: 'rgba(255,255,255,0.07)' }}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/home" className="flex items-center gap-2" style={{ textDecoration: 'none' }}>
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: 'var(--brand-amber)' }}>
              <Ticket className="w-4 h-4" style={{ color: '#0f1729' }} />
            </div>
            <span className="text-xl font-bold tracking-tight" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
              uduX<span style={{ color: 'var(--brand-amber)' }}>Pass</span>
            </span>
          </Link>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-8">
            {[{ path: '/home', label: 'Home' }, { path: '/events', label: 'Events' }].map(link => (
              <Link key={link.path} to={link.path} className="text-sm font-medium transition-colors"
                style={{ color: isActive(link.path) ? 'var(--brand-amber)' : '#94a3b8', fontFamily: 'var(--font-body)' }}>
                {link.label}
              </Link>
            ))}
          </div>

          {/* Right Actions */}
          <div className="flex items-center gap-3">
            <button onClick={toggleCart} className="relative p-2 rounded-lg" style={{ color: '#94a3b8' }} aria-label="Cart">
              <ShoppingCart className="w-5 h-5" />
              {cartCount > 0 && (
                <Badge className="absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs"
                  style={{ background: 'var(--brand-amber)', color: '#0f1729', fontWeight: 700 }}>
                  {cartCount}
                </Badge>
              )}
            </button>

            {isAuthenticated && user ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button className="flex items-center gap-2 px-3 py-1.5 rounded-lg"
                    style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }}>
                    <div className="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold"
                      style={{ background: 'var(--brand-amber)', color: '#0f1729' }}>
                      {user.firstName?.[0]?.toUpperCase() || 'U'}
                    </div>
                    <span className="text-sm font-medium hidden sm:block">{user.firstName}</span>
                    <ChevronDown className="w-3 h-3 opacity-60" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-52"
                  style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }}>
                  <DropdownMenuLabel className="text-xs opacity-60">{user.email}</DropdownMenuLabel>
                  <DropdownMenuSeparator style={{ background: 'rgba(255,255,255,0.08)' }} />
                  <DropdownMenuItem onClick={() => navigate('/profile')} className="cursor-pointer">
                    <User className="mr-2 h-4 w-4" /> My Profile
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => navigate('/profile?tab=orders')} className="cursor-pointer">
                    <Ticket className="mr-2 h-4 w-4" /> My Tickets
                  </DropdownMenuItem>
                  <DropdownMenuSeparator style={{ background: 'rgba(255,255,255,0.08)' }} />
                  <DropdownMenuItem onClick={handleLogout} className="cursor-pointer" style={{ color: '#f87171' }}>
                    <LogOut className="mr-2 h-4 w-4" /> Sign Out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <div className="flex items-center gap-2">
                <Button variant="ghost" size="sm" onClick={() => navigate('/login')} style={{ color: '#94a3b8' }}>Sign In</Button>
                <Button size="sm" onClick={() => navigate('/register')}
                  style={{ background: 'var(--brand-amber)', color: '#0f1729', fontWeight: 700, fontFamily: 'var(--font-display)' }}>
                  Get Tickets
                </Button>
              </div>
            )}

            <Button variant="ghost" size="sm" className="md:hidden" onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)} style={{ color: '#94a3b8' }}>
              {isMobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </Button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMobileMenuOpen && (
          <div className="md:hidden py-4 border-t" style={{ borderColor: 'rgba(255,255,255,0.07)' }}>
            <div className="flex flex-col gap-2">
              {[{ path: '/home', label: 'Home' }, { path: '/events', label: 'Events' }].map(link => (
                <Link key={link.path} to={link.path} className="px-3 py-2 rounded-lg text-sm font-medium"
                  style={{ color: isActive(link.path) ? 'var(--brand-amber)' : '#94a3b8', background: isActive(link.path) ? 'rgba(245,158,11,0.1)' : 'transparent' }}
                  onClick={() => setIsMobileMenuOpen(false)}>
                  {link.label}
                </Link>
              ))}
              {!isAuthenticated && (
                <div className="flex flex-col gap-2 pt-3 border-t" style={{ borderColor: 'rgba(255,255,255,0.07)' }}>
                  <Button variant="ghost" className="justify-start" style={{ color: '#94a3b8' }} onClick={() => { navigate('/login'); setIsMobileMenuOpen(false) }}>Sign In</Button>
                  <Button className="justify-start" style={{ background: 'var(--brand-amber)', color: '#0f1729', fontWeight: 700 }} onClick={() => { navigate('/register'); setIsMobileMenuOpen(false) }}>Get Tickets</Button>
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

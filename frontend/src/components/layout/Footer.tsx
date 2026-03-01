/*
 * Footer — uduXPass Design System
 * Dark navy/amber, Syne headings
 */
import React from 'react'
import { Link } from 'react-router-dom'
import { Ticket, Instagram, Twitter, Facebook, Youtube } from 'lucide-react'

const Footer: React.FC = () => {
  const year = new Date().getFullYear()

  return (
    <footer style={{ background: 'var(--brand-surface)', borderTop: '1px solid rgba(255,255,255,0.07)' }}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12 mb-12">
          {/* Brand */}
          <div className="md:col-span-1">
            <div className="flex items-center gap-2 mb-4">
              <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: 'var(--brand-amber)' }}>
                <Ticket className="w-4 h-4" style={{ color: '#0f1729' }} />
              </div>
              <span className="text-lg font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
                uduX<span style={{ color: 'var(--brand-amber)' }}>Pass</span>
              </span>
            </div>
            <p className="text-sm mb-6" style={{ color: '#64748b', lineHeight: 1.7 }}>
              Nigeria's premier ticketing platform. Discover, book, and experience the best events across the country.
            </p>
            <div className="flex gap-3">
              {[
                { icon: Instagram, href: '#' },
                { icon: Twitter, href: '#' },
                { icon: Facebook, href: '#' },
                { icon: Youtube, href: '#' },
              ].map(({ icon: Icon, href }) => (
                <a key={href} href={href}
                  className="w-8 h-8 rounded-lg flex items-center justify-center transition-colors"
                  style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.08)', color: '#64748b' }}>
                  <Icon className="w-4 h-4" />
                </a>
              ))}
            </div>
          </div>

          {/* Links */}
          {[
            {
              title: 'Platform',
              links: [
                { label: 'Browse Events', to: '/events' },
                { label: 'My Tickets', to: '/profile' },
                { label: 'Create Account', to: '/register' },
                { label: 'Sign In', to: '/login' },
              ]
            },
            {
              title: 'Company',
              links: [
                { label: 'About Us', to: '/about' },
                { label: 'Contact', to: '/contact' },
                { label: 'Careers', to: '/careers' },
                { label: 'Blog', to: '/blog' },
              ]
            },
            {
              title: 'Legal',
              links: [
                { label: 'Terms of Service', to: '/terms' },
                { label: 'Privacy Policy', to: '/privacy' },
                { label: 'Cookie Policy', to: '/cookies' },
                { label: 'Refund Policy', to: '/refunds' },
              ]
            },
          ].map(section => (
            <div key={section.title}>
              <h4 className="text-sm font-bold mb-4 tracking-wide" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
                {section.title}
              </h4>
              <ul className="space-y-3">
                {section.links.map(link => (
                  <li key={link.label}>
                    <Link to={link.to} className="text-sm transition-colors"
                      style={{ color: '#64748b', textDecoration: 'none' }}
                      onMouseEnter={e => (e.target as HTMLElement).style.color = 'var(--brand-amber)'}
                      onMouseLeave={e => (e.target as HTMLElement).style.color = '#64748b'}>
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Bottom */}
        <div className="flex flex-col md:flex-row items-center justify-between gap-4 pt-8"
          style={{ borderTop: '1px solid rgba(255,255,255,0.07)' }}>
          <p className="text-xs" style={{ color: '#475569' }}>
            © {year} uduXPass. All rights reserved.
          </p>
          <p className="text-xs" style={{ color: '#334155' }}>
            Made with ♥ in Nigeria
          </p>
        </div>
      </div>
    </footer>
  )
}

export default Footer

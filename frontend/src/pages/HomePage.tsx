/*
 * HomePage — uduXPass Design System
 * Dark hero with atmospheric concert imagery, amber CTAs, category browsing
 */
import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { eventsAPI } from '@/services/api'
import { Event } from '@/types/api'
import { Calendar, MapPin, ArrowRight, Ticket, Music, Mic2, Laugh, Trophy, Zap, Users, Star, TrendingUp } from 'lucide-react'

const formatDate = (dateStr: string) => {
  try {
    return new Date(dateStr).toLocaleDateString('en-NG', { day: 'numeric', month: 'short', year: 'numeric' })
  } catch { return dateStr }
}

const formatCurrency = (amount: number) => {
  if (!amount || isNaN(amount)) return 'Free'
  return new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN', minimumFractionDigits: 0 }).format(amount)
}

const getLowestPrice = (event: Event): number => {
  if (!event.ticket_tiers?.length) return 0
  const prices = event.ticket_tiers.map(t => t.price).filter(p => p > 0)
  return prices.length ? Math.min(...prices) : 0
}

const CATEGORIES = [
  { label: 'Concerts', icon: Music, color: '#F59E0B' },
  { label: 'Comedy', icon: Laugh, color: '#10b981' },
  { label: 'Festivals', icon: Zap, color: '#8b5cf6' },
  { label: 'Sports', icon: Trophy, color: '#ef4444' },
  { label: 'Talks', icon: Mic2, color: '#3b82f6' },
]

const FALLBACK_IMAGES = [
  'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=600&q=80',
  'https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?w=600&q=80',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=600&q=80',
  'https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=600&q=80',
  'https://images.unsplash.com/photo-1506157786151-b8491531f063?w=600&q=80',
  'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=600&q=80',
]

const HomePage: React.FC = () => {
  const [featuredEvents, setFeaturedEvents] = useState<Event[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => { loadFeaturedEvents() }, [])

  const loadFeaturedEvents = async () => {
    try {
      setIsLoading(true)
      setError(null)
      const response = await eventsAPI.getEvents({ limit: 6, status: 'published' })
      if (response.success && response.data) {
        const data = response.data as any
        setFeaturedEvents(Array.isArray(data) ? data : data.events || data.data || [])
      } else {
        setError('Could not load events')
      }
    } catch {
      setError('Could not load events')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div style={{ background: 'var(--brand-navy)', minHeight: '100vh' }}>

      {/* Hero */}
      <section className="relative overflow-hidden" style={{ minHeight: '88vh', display: 'flex', alignItems: 'center' }}>
        <div className="absolute inset-0">
          <img
            src="https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=1600&q=80"
            alt="Concert crowd"
            className="w-full h-full object-cover"
            style={{ opacity: 0.35 }}
          />
          <div className="absolute inset-0" style={{ background: 'linear-gradient(135deg, rgba(15,23,41,0.95) 0%, rgba(15,23,41,0.7) 50%, rgba(15,23,41,0.9) 100%)' }} />
        </div>
        <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="max-w-3xl">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full mb-6 animate-fade-up"
              style={{ background: 'rgba(245,158,11,0.15)', border: '1px solid rgba(245,158,11,0.3)' }}>
              <div className="w-2 h-2 rounded-full animate-pulse-amber" style={{ background: 'var(--brand-amber)' }} />
              <span className="text-xs font-semibold tracking-widest uppercase" style={{ color: 'var(--brand-amber)', fontFamily: 'var(--font-display)' }}>
                Nigeria's Premier Ticketing Platform
              </span>
            </div>
            <h1 className="animate-fade-up-delay-1 mb-6"
              style={{ fontSize: 'clamp(2.5rem, 6vw, 4.5rem)', fontFamily: 'var(--font-display)', color: '#f1f5f9', lineHeight: 1.05 }}>
              Your Ticket to<br />
              <span style={{ color: 'var(--brand-amber)' }}>Unforgettable</span><br />
              Experiences
            </h1>
            <p className="animate-fade-up-delay-2 mb-10 text-lg" style={{ color: '#94a3b8', maxWidth: '520px', lineHeight: 1.7 }}>
              Discover the hottest concerts, festivals, and live events across Nigeria.
              Secure your spot instantly with MoMo or card payments.
            </p>
            <div className="flex flex-wrap gap-4 animate-fade-up-delay-3">
              <Link to="/events">
                <Button size="lg" className="gap-2 px-8 py-6 text-base font-bold rounded-xl"
                  style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)', boxShadow: '0 4px 20px rgba(245,158,11,0.4)' }}>
                  Browse Events <ArrowRight className="w-5 h-5" />
                </Button>
              </Link>
              <Link to="/register">
                <Button size="lg" variant="outline" className="gap-2 px-8 py-6 text-base font-semibold rounded-xl"
                  style={{ borderColor: 'rgba(255,255,255,0.2)', color: '#f1f5f9', background: 'rgba(255,255,255,0.05)' }}>
                  Create Account
                </Button>
              </Link>
            </div>
            <div className="flex flex-wrap gap-8 mt-14 animate-fade-up-delay-3">
              {[{ value: '50K+', label: 'Tickets Sold' }, { value: '200+', label: 'Events Hosted' }, { value: '99.9%', label: 'Uptime' }].map(stat => (
                <div key={stat.label}>
                  <div className="text-2xl font-bold" style={{ fontFamily: 'var(--font-display)', color: 'var(--brand-amber)' }}>{stat.value}</div>
                  <div className="text-sm" style={{ color: '#64748b' }}>{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Categories */}
      <section className="py-16" style={{ background: 'var(--brand-surface)', borderTop: '1px solid rgba(255,255,255,0.06)', borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-xl font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Browse by Category</h2>
            <Link to="/events" className="text-sm font-medium flex items-center gap-1" style={{ color: 'var(--brand-amber)' }}>
              View all <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
          <div className="flex flex-wrap gap-3">
            {CATEGORIES.map(cat => {
              const Icon = cat.icon
              return (
                <Link key={cat.label} to={`/events?category=${cat.label.toLowerCase()}`} style={{ textDecoration: 'none' }}>
                  <div className="flex items-center gap-2 px-5 py-3 rounded-full cursor-pointer transition-all duration-200 hover:-translate-y-0.5"
                    style={{ background: 'rgba(255,255,255,0.05)', border: '1px solid rgba(255,255,255,0.1)', color: '#f1f5f9' }}>
                    <Icon className="w-4 h-4" style={{ color: cat.color }} />
                    <span className="text-sm font-medium">{cat.label}</span>
                  </div>
                </Link>
              )
            })}
          </div>
        </div>
      </section>

      {/* Featured Events */}
      <section className="py-20" style={{ background: 'var(--brand-navy)' }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-end justify-between mb-12">
            <div>
              <p className="text-xs font-semibold tracking-widest uppercase mb-2" style={{ color: 'var(--brand-amber)', fontFamily: 'var(--font-display)' }}>On Sale Now</p>
              <h2 className="text-3xl md:text-4xl font-bold" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Featured Events</h2>
            </div>
            <Link to="/events">
              <Button variant="outline" size="sm" className="gap-2"
                style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8', background: 'transparent' }}>
                View All <ArrowRight className="w-4 h-4" />
              </Button>
            </Link>
          </div>
          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="rounded-2xl overflow-hidden animate-pulse" style={{ background: 'var(--brand-surface)', height: '380px' }} />
              ))}
            </div>
          ) : error || featuredEvents.length === 0 ? (
            <div className="text-center py-20">
              <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4"
                style={{ background: 'rgba(245,158,11,0.1)', border: '1px solid rgba(245,158,11,0.2)' }}>
                <Ticket className="w-8 h-8" style={{ color: 'var(--brand-amber)' }} />
              </div>
              <h3 className="text-xl font-bold mb-2" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>No Events Available Yet</h3>
              <p className="mb-6" style={{ color: '#64748b' }}>New events are added regularly. Check back soon!</p>
              <Button onClick={loadFeaturedEvents} variant="outline" style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8' }}>Refresh</Button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {featuredEvents.map((event, index) => (
                <Link key={event.id} to={`/events/${event.id}`} style={{ textDecoration: 'none' }}>
                  <div className="rounded-2xl overflow-hidden cursor-pointer group transition-all duration-300 hover:-translate-y-1"
                    style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.07)', boxShadow: '0 4px 24px rgba(0,0,0,0.3)' }}>
                    <div className="relative h-52 overflow-hidden">
                      <img
                        src={event.event_image_url || event.banner_image_url || FALLBACK_IMAGES[index % FALLBACK_IMAGES.length]}
                        alt={event.name}
                        className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                        onError={(e) => { (e.target as HTMLImageElement).src = FALLBACK_IMAGES[index % FALLBACK_IMAGES.length] }}
                      />
                      <div className="absolute inset-0" style={{ background: 'linear-gradient(to top, rgba(13,21,38,0.8) 0%, transparent 60%)' }} />
                      <div className="absolute top-3 right-3">
                        <Badge className="text-xs font-bold px-2 py-1" style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)' }}>On Sale</Badge>
                      </div>
                    </div>
                    <div className="p-5">
                      <h3 className="font-bold text-base mb-3 line-clamp-2" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>{event.name}</h3>
                      <div className="flex flex-col gap-1.5 mb-4">
                        <div className="flex items-center gap-2 text-sm" style={{ color: '#64748b' }}>
                          <Calendar className="w-3.5 h-3.5 flex-shrink-0" />
                          <span>{formatDate(event.event_date)}</span>
                        </div>
                        <div className="flex items-center gap-2 text-sm" style={{ color: '#64748b' }}>
                          <MapPin className="w-3.5 h-3.5 flex-shrink-0" />
                          <span className="line-clamp-1">{event.venue_name || (event.venue as any)?.name || 'Venue TBA'}</span>
                        </div>
                      </div>
                      <div className="flex items-center justify-between pt-3" style={{ borderTop: '1px solid rgba(255,255,255,0.07)' }}>
                        <div>
                          <p className="text-xs mb-0.5" style={{ color: '#64748b' }}>From</p>
                          <p className="text-lg font-bold" style={{ fontFamily: 'var(--font-display)', color: 'var(--brand-amber)' }}>{formatCurrency(getLowestPrice(event))}</p>
                        </div>
                        <div className="flex items-center gap-1 text-xs px-3 py-1.5 rounded-full"
                          style={{ background: 'rgba(245,158,11,0.1)', color: 'var(--brand-amber)', border: '1px solid rgba(245,158,11,0.2)' }}>
                          <Ticket className="w-3 h-3" /> Get Tickets
                        </div>
                      </div>
                    </div>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      </section>

      {/* Trust Bar */}
      <section className="py-16" style={{ background: 'var(--brand-surface)', borderTop: '1px solid rgba(255,255,255,0.06)' }}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
            {[
              { icon: Ticket, label: 'Instant E-Tickets', desc: 'Delivered to your email & app' },
              { icon: Users, label: '50,000+ Users', desc: 'Trust uduXPass for events' },
              { icon: Star, label: 'Verified Events', desc: 'All events are vetted' },
              { icon: TrendingUp, label: 'MoMo & Card', desc: 'Flexible payment options' },
            ].map(item => {
              const Icon = item.icon
              return (
                <div key={item.label} className="flex flex-col items-center gap-3">
                  <div className="w-12 h-12 rounded-xl flex items-center justify-center"
                    style={{ background: 'rgba(245,158,11,0.1)', border: '1px solid rgba(245,158,11,0.2)' }}>
                    <Icon className="w-6 h-6" style={{ color: 'var(--brand-amber)' }} />
                  </div>
                  <div>
                    <p className="font-bold text-sm" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>{item.label}</p>
                    <p className="text-xs mt-0.5" style={{ color: '#64748b' }}>{item.desc}</p>
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="relative overflow-hidden py-24">
        <div className="absolute inset-0">
          <img src="https://images.unsplash.com/photo-1506157786151-b8491531f063?w=1600&q=80" alt="Festival" className="w-full h-full object-cover" style={{ opacity: 0.2 }} />
          <div className="absolute inset-0" style={{ background: 'linear-gradient(135deg, rgba(15,23,41,0.97) 0%, rgba(30,45,74,0.9) 100%)' }} />
        </div>
        <div className="relative z-10 max-w-3xl mx-auto px-4 text-center">
          <h2 className="text-3xl md:text-5xl font-bold mb-6" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
            Ready to Experience<br /><span style={{ color: 'var(--brand-amber)' }}>Something Legendary?</span>
          </h2>
          <p className="text-lg mb-10" style={{ color: '#94a3b8' }}>
            Join thousands of Nigerians who use uduXPass to discover and attend the best events.
          </p>
          <div className="flex flex-wrap gap-4 justify-center">
            <Link to="/events">
              <Button size="lg" className="px-10 py-6 text-base font-bold rounded-xl gap-2"
                style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)', boxShadow: '0 4px 20px rgba(245,158,11,0.4)' }}>
                Browse Events <ArrowRight className="w-5 h-5" />
              </Button>
            </Link>
            <Link to="/register">
              <Button size="lg" variant="outline" className="px-10 py-6 text-base font-semibold rounded-xl"
                style={{ borderColor: 'rgba(255,255,255,0.2)', color: '#f1f5f9', background: 'rgba(255,255,255,0.05)' }}>
                Create Free Account
              </Button>
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}

export default HomePage

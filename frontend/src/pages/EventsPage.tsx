/*
 * EventsPage — uduXPass Design System
 * Dark navy/amber, Syne headings — premium event discovery
 */
import React, { useState, useEffect, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { eventsAPI } from '../services/api'
import { Event, PaginatedResponse, ApiResponse } from '../types/api'
import { Calendar, MapPin, Search, Ticket, ArrowRight, SlidersHorizontal, X } from 'lucide-react'

const formatDate = (d: string) => {
  try { return new Date(d).toLocaleDateString('en-NG', { day: 'numeric', month: 'short', year: 'numeric' }) }
  catch { return d }
}
const formatCurrency = (n: number) => {
  if (!n || isNaN(n)) return 'Free'
  return new Intl.NumberFormat('en-NG', { style: 'currency', currency: 'NGN', minimumFractionDigits: 0 }).format(n)
}
const getLowestPrice = (event: Event): number => {
  if (!event.ticket_tiers?.length) return 0
  const prices = event.ticket_tiers.map(t => t.price).filter(p => p > 0)
  return prices.length ? Math.min(...prices) : 0
}
const getAvailableTickets = (event: Event): number => {
  if (!event.ticket_tiers?.length) return 0
  return event.ticket_tiers.reduce((sum, t) => sum + (t.available_quantity ?? t.quantity ?? 0), 0)
}

const FALLBACK_IMAGES = [
  'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=600&q=80',
  'https://images.unsplash.com/photo-1501281668745-f7f57925c3b4?w=600&q=80',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=600&q=80',
  'https://images.unsplash.com/photo-1514525253161-7a46d19cd819?w=600&q=80',
  'https://images.unsplash.com/photo-1506157786151-b8491531f063?w=600&q=80',
  'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=600&q=80',
]

const SORT_OPTIONS = [
  { value: 'date', label: 'Date' },
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
]

const EventsPage: React.FC = () => {
  const [events, setEvents] = useState<Event[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [city, setCity] = useState('')
  const [sort, setSort] = useState('date')
  const [showFilters, setShowFilters] = useState(false)
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [total, setTotal] = useState(0)

  const loadEvents = useCallback(async (p = 1, s = search, c = city) => {
    try {
      setIsLoading(true)
      setError(null)
      const response: ApiResponse<PaginatedResponse<Event>> = await eventsAPI.getEvents({
        page: p, limit: 12, search: s || undefined, city: c || undefined, status: 'published'
      })
      if (response.success && response.data) {
        const data = response.data as any
        const arr: Event[] = Array.isArray(data) ? data : data.events || data.data || []
        setEvents(arr)
        setTotalPages(data.totalPages || data.total_pages || 1)
        setTotal(data.total || arr.length)
      } else {
        setError('Could not load events')
      }
    } catch {
      setError('Could not load events')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => { loadEvents(1) }, [])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setPage(1)
    loadEvents(1, search, city)
  }

  const clearFilters = () => {
    setSearch('')
    setCity('')
    setSort('date')
    setPage(1)
    loadEvents(1, '', '')
  }

  const hasFilters = search || city || sort !== 'date'

  const sortedEvents = [...events].sort((a, b) => {
    if (sort === 'price_asc') return getLowestPrice(a) - getLowestPrice(b)
    if (sort === 'price_desc') return getLowestPrice(b) - getLowestPrice(a)
    return new Date(a.event_date).getTime() - new Date(b.event_date).getTime()
  })

  return (
    <div style={{ background: 'var(--brand-navy)', minHeight: '100vh' }}>
      {/* Page Header */}
      <div className="relative overflow-hidden py-16"
        style={{ background: 'var(--brand-surface)', borderBottom: '1px solid rgba(255,255,255,0.07)' }}>
        <div className="absolute inset-0 opacity-10"
          style={{ backgroundImage: 'radial-gradient(ellipse at 70% 50%, var(--brand-amber) 0%, transparent 60%)' }} />
        <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <p className="text-xs font-semibold tracking-widest uppercase mb-2" style={{ color: 'var(--brand-amber)', fontFamily: 'var(--font-display)' }}>
            Discover
          </p>
          <h1 className="text-3xl md:text-5xl font-bold mb-4" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>
            All Events
          </h1>
          <p className="text-base mb-8" style={{ color: '#64748b' }}>
            {total > 0 ? `${total} events available` : 'Find your next unforgettable experience'}
          </p>
          {/* Search Bar */}
          <form onSubmit={handleSearch} className="flex gap-3 max-w-2xl">
            <div className="relative flex-1">
              <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4" style={{ color: '#475569' }} />
              <Input
                value={search}
                onChange={e => setSearch(e.target.value)}
                placeholder="Search events, artists, venues..."
                className="pl-11 h-12 text-sm"
                style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', color: '#f1f5f9' }}
              />
            </div>
            <Button type="submit" className="h-12 px-6 font-bold"
              style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)' }}>
              Search
            </Button>
            <Button type="button" variant="outline" className="h-12 px-4" onClick={() => setShowFilters(!showFilters)}
              style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8', background: 'transparent' }}>
              <SlidersHorizontal className="w-4 h-4" />
            </Button>
          </form>

          {/* Filter Row */}
          {showFilters && (
            <div className="flex flex-wrap gap-3 mt-4 items-center animate-fade-up">
              <Input
                value={city}
                onChange={e => setCity(e.target.value)}
                placeholder="Filter by city..."
                className="h-9 text-sm w-48"
                style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', color: '#f1f5f9' }}
              />
              <select
                value={sort}
                onChange={e => setSort(e.target.value)}
                className="h-9 px-3 text-sm rounded-md"
                style={{ background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', color: '#f1f5f9' }}
              >
                {SORT_OPTIONS.map(o => <option key={o.value} value={o.value} style={{ background: '#0d1526' }}>{o.label}</option>)}
              </select>
              {hasFilters && (
                <button onClick={clearFilters} className="flex items-center gap-1 text-xs px-3 py-1.5 rounded-full"
                  style={{ color: '#f87171', background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)' }}>
                  <X className="w-3 h-3" /> Clear
                </button>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Events Grid */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(9)].map((_, i) => (
              <div key={i} className="rounded-2xl overflow-hidden animate-pulse" style={{ background: 'var(--brand-surface)', height: '380px' }} />
            ))}
          </div>
        ) : error ? (
          <div className="text-center py-20">
            <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4"
              style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.2)' }}>
              <Ticket className="w-8 h-8" style={{ color: '#f87171' }} />
            </div>
            <h3 className="text-xl font-bold mb-2" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>Could Not Load Events</h3>
            <p className="mb-6" style={{ color: '#64748b' }}>There was a problem fetching events. Please try again.</p>
            <Button onClick={() => loadEvents(1)} style={{ background: 'var(--brand-amber)', color: '#0f1729', fontWeight: 700 }}>Retry</Button>
          </div>
        ) : sortedEvents.length === 0 ? (
          <div className="text-center py-20">
            <div className="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4"
              style={{ background: 'rgba(245,158,11,0.1)', border: '1px solid rgba(245,158,11,0.2)' }}>
              <Ticket className="w-8 h-8" style={{ color: 'var(--brand-amber)' }} />
            </div>
            <h3 className="text-xl font-bold mb-2" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>No Events Found</h3>
            <p className="mb-6" style={{ color: '#64748b' }}>
              {hasFilters ? 'No events match your filters. Try adjusting your search.' : 'No events are available right now. Check back soon!'}
            </p>
            {hasFilters && <Button onClick={clearFilters} variant="outline" style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8' }}>Clear Filters</Button>}
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
              {sortedEvents.map((event, index) => (
                <Link key={event.id} to={`/events/${event.id}`} style={{ textDecoration: 'none' }}>
                  <div className="rounded-2xl overflow-hidden cursor-pointer group transition-all duration-300 hover:-translate-y-1 h-full flex flex-col"
                    style={{ background: 'var(--brand-surface)', border: '1px solid rgba(255,255,255,0.07)', boxShadow: '0 4px 24px rgba(0,0,0,0.3)' }}>
                    <div className="relative h-52 overflow-hidden flex-shrink-0">
                      <img
                        src={event.event_image_url || event.banner_image_url || FALLBACK_IMAGES[index % FALLBACK_IMAGES.length]}
                        alt={event.name}
                        className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                        onError={(e) => { (e.target as HTMLImageElement).src = FALLBACK_IMAGES[index % FALLBACK_IMAGES.length] }}
                      />
                      <div className="absolute inset-0" style={{ background: 'linear-gradient(to top, rgba(13,21,38,0.8) 0%, transparent 60%)' }} />
                      <div className="absolute top-3 right-3">
                        <Badge className="text-xs font-bold px-2 py-1"
                          style={{ background: 'var(--brand-amber)', color: '#0f1729', fontFamily: 'var(--font-display)' }}>
                          {event.status === 'published' ? 'On Sale' : event.status}
                        </Badge>
                      </div>
                    </div>
                    <div className="p-5 flex flex-col flex-1">
                      <h3 className="font-bold text-base mb-3 line-clamp-2" style={{ fontFamily: 'var(--font-display)', color: '#f1f5f9' }}>{event.name}</h3>
                      <div className="flex flex-col gap-1.5 mb-4 flex-1">
                        <div className="flex items-center gap-2 text-sm" style={{ color: '#64748b' }}>
                          <Calendar className="w-3.5 h-3.5 flex-shrink-0" />
                          <span>{formatDate(event.event_date)}</span>
                        </div>
                        <div className="flex items-center gap-2 text-sm" style={{ color: '#64748b' }}>
                          <MapPin className="w-3.5 h-3.5 flex-shrink-0" />
                          <span className="line-clamp-1">
                            {event.venue_name || (event.venue as any)?.name || 'Venue TBA'}
                            {(event.venue as any)?.city ? `, ${(event.venue as any).city}` : ''}
                          </span>
                        </div>
                      </div>
                      <div className="flex items-center justify-between pt-3" style={{ borderTop: '1px solid rgba(255,255,255,0.07)' }}>
                        <div>
                          <p className="text-xs mb-0.5" style={{ color: '#64748b' }}>From</p>
                          <p className="text-lg font-bold" style={{ fontFamily: 'var(--font-display)', color: 'var(--brand-amber)' }}>
                            {formatCurrency(getLowestPrice(event))}
                          </p>
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

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex items-center justify-center gap-2">
                <Button variant="outline" onClick={() => { setPage(p => p - 1); loadEvents(page - 1) }} disabled={page === 1}
                  style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8', background: 'transparent' }}>
                  Previous
                </Button>
                {Array.from({ length: Math.min(totalPages, 7) }, (_, i) => i + 1).map(p => (
                  <Button key={p} variant={p === page ? 'default' : 'outline'}
                    onClick={() => { setPage(p); loadEvents(p) }}
                    className="min-w-[40px]"
                    style={p === page
                      ? { background: 'var(--brand-amber)', color: '#0f1729', fontWeight: 700 }
                      : { borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8', background: 'transparent' }}>
                    {p}
                  </Button>
                ))}
                <Button variant="outline" onClick={() => { setPage(p => p + 1); loadEvents(page + 1) }} disabled={page === totalPages}
                  style={{ borderColor: 'rgba(255,255,255,0.15)', color: '#94a3b8', background: 'transparent' }}>
                  Next
                </Button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}

export default EventsPage

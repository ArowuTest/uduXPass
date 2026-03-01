/*
 * EventsPage — uduXPass Design System
 * Design: Dark navy (#0f1729) background, amber (#F59E0B) accents
 * Typography: Playfair Display (headings) + Inter (body)
 */
import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { eventsAPI } from '../services/api';
import { Event } from '../types/api';

const CATEGORIES = [
  { label: 'All', value: '' },
  { label: 'Concerts', value: 'concerts' },
  { label: 'Festivals', value: 'festivals' },
  { label: 'Comedy', value: 'comedy' },
  { label: 'Sports', value: 'sports' },
  { label: 'Theatre', value: 'theatre' },
  { label: 'Tech', value: 'tech' },
];

function formatDate(dateStr: string) {
  if (!dateStr) return 'TBA';
  try {
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-NG', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' });
  } catch { return dateStr; }
}

function formatPrice(tiers: any[]) {
  if (!tiers || tiers.length === 0) return 'Free';
  const prices = tiers.map((t: any) => Number(t.price || 0)).filter(p => p > 0);
  if (prices.length === 0) return 'Free';
  const min = Math.min(...prices);
  const max = Math.max(...prices);
  const fmt = (n: number) => `\u20A6${n.toLocaleString('en-NG')}`;
  return min === max ? fmt(min) : `${fmt(min)} – ${fmt(max)}`;
}

export default function EventsPage() {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [search, setSearch] = useState('');
  const [city, setCity] = useState('');
  const [activeCategory, setActiveCategory] = useState('');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);

  const fetchEvents = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await eventsAPI.getEvents({ page, limit: 12, search, city });
      if (res.success && res.data) {
        // api.ts transforms backend response: {data: events[], pagination: {}}
        const raw = res.data as any;
        // Handle multiple possible shapes from the transformer
        let evtData: any[] = [];
        let pagination: any = null;
        if (Array.isArray(raw.data)) {
          evtData = raw.data;
          pagination = raw.pagination;
        } else if (Array.isArray(raw)) {
          evtData = raw;
        } else if (raw.events) {
          evtData = raw.events;
          pagination = raw.pagination;
        }
        setEvents(evtData);
        if (pagination) {
          setTotalPages(pagination.total_pages || 1);
          setTotalCount(pagination.total || evtData.length);
        } else {
          setTotalPages(1);
          setTotalCount(evtData.length);
        }
      } else {
        setError(res.error || 'Failed to load events');
      }
    } catch (err) {
      setError('Unable to connect. Please try again.');
    } finally {
      setLoading(false);
    }
  }, [page, search, city]);

  useEffect(() => { fetchEvents(); }, [fetchEvents]);

  const handleSearch = (e: React.FormEvent) => { e.preventDefault(); setPage(1); };

  const getEventImage = (event: any) =>
    event.event_image_url || event.eventImageUrl || event.banner_image_url || null;

  const getEventDate = (event: any) =>
    event.event_date || event.eventDate || event.date || '';

  const getVenueInfo = (event: any) => {
    if (event.venue) return { name: event.venue.name, city: event.venue.city };
    return { name: event.venue_name || event.venueName || 'Venue TBA', city: event.venue_city || event.venueCity || '' };
  };

  const getTiers = (event: any) => event.ticket_tiers || event.ticketTiers || [];

  return (
    <div style={{ background: 'var(--color-navy, #0f1729)', minHeight: '100vh', color: '#fff' }}>
      {/* Hero search */}
      <div style={{
        background: 'linear-gradient(135deg, #0f1729 0%, #1a2744 50%, #0f1729 100%)',
        borderBottom: '1px solid rgba(245,158,11,0.15)',
        padding: '3rem 1rem 2.5rem',
      }}>
        <div style={{ maxWidth: '900px', margin: '0 auto', textAlign: 'center' }}>
          <h1 style={{
            fontFamily: "'Playfair Display', serif",
            fontSize: 'clamp(2rem, 5vw, 3rem)',
            fontWeight: 700,
            color: '#fff',
            marginBottom: '0.5rem',
          }}>Discover Amazing Events</h1>
          <p style={{ color: 'rgba(255,255,255,0.6)', marginBottom: '2rem', fontSize: '1.05rem' }}>
            Find and book tickets to the hottest concerts, festivals, and shows
          </p>
          <form onSubmit={handleSearch} style={{ display: 'flex', gap: '0.75rem', maxWidth: '700px', margin: '0 auto', flexWrap: 'wrap' }}>
            <div style={{ flex: 1, minWidth: '200px', position: 'relative' }}>
              <input
                type="text"
                placeholder="Search events, artists, or venues..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                style={{
                  width: '100%', padding: '0.875rem 1rem 0.875rem 1rem',
                  background: 'rgba(255,255,255,0.08)', border: '1px solid rgba(255,255,255,0.15)',
                  borderRadius: '0.625rem', color: '#fff', fontSize: '0.95rem', outline: 'none',
                  boxSizing: 'border-box',
                }}
              />
            </div>
            <select
              value={city}
              onChange={e => { setCity(e.target.value); setPage(1); }}
              style={{
                padding: '0.875rem 1rem', background: 'rgba(255,255,255,0.08)',
                border: '1px solid rgba(255,255,255,0.15)', borderRadius: '0.625rem',
                color: '#fff', fontSize: '0.95rem', cursor: 'pointer', minWidth: '130px',
              }}
            >
              <option value="" style={{ background: '#1a2744' }}>All Cities</option>
              <option value="Lagos" style={{ background: '#1a2744' }}>Lagos</option>
              <option value="Abuja" style={{ background: '#1a2744' }}>Abuja</option>
              <option value="Port Harcourt" style={{ background: '#1a2744' }}>Port Harcourt</option>
            </select>
            <button type="submit" style={{
              padding: '0.875rem 1.5rem', background: '#F59E0B', color: '#000',
              border: 'none', borderRadius: '0.625rem', fontWeight: 700, fontSize: '0.95rem', cursor: 'pointer',
            }}>Search</button>
          </form>
        </div>
      </div>

      {/* Category pills */}
      <div style={{ borderBottom: '1px solid rgba(255,255,255,0.07)', padding: '1rem 1.5rem' }}>
        <div style={{ maxWidth: '1200px', margin: '0 auto', display: 'flex', gap: '0.5rem', overflowX: 'auto', paddingBottom: '0.25rem' }}>
          {CATEGORIES.map(cat => (
            <button key={cat.value} onClick={() => { setActiveCategory(cat.value); setPage(1); }}
              style={{
                padding: '0.4rem 1.1rem', borderRadius: '999px',
                border: activeCategory === cat.value ? '1px solid #F59E0B' : '1px solid rgba(255,255,255,0.15)',
                background: activeCategory === cat.value ? 'rgba(245,158,11,0.15)' : 'transparent',
                color: activeCategory === cat.value ? '#F59E0B' : 'rgba(255,255,255,0.6)',
                fontSize: '0.875rem', fontWeight: activeCategory === cat.value ? 600 : 400,
                cursor: 'pointer', whiteSpace: 'nowrap',
              }}>{cat.label}</button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem 1.5rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
          <p style={{ color: 'rgba(255,255,255,0.5)', fontSize: '0.9rem' }}>
            {loading ? 'Loading events...' : `${totalCount} event${totalCount !== 1 ? 's' : ''} found`}
          </p>
        </div>

        {loading && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
            {[1,2,3,4,5,6].map(i => (
              <div key={i} style={{ background: 'rgba(255,255,255,0.04)', borderRadius: '1rem', overflow: 'hidden', border: '1px solid rgba(255,255,255,0.07)' }}>
                <div style={{ height: '200px', background: 'rgba(255,255,255,0.06)' }} />
                <div style={{ padding: '1.25rem' }}>
                  <div style={{ height: '1rem', background: 'rgba(255,255,255,0.06)', borderRadius: '0.25rem', marginBottom: '0.75rem', width: '70%' }} />
                  <div style={{ height: '0.75rem', background: 'rgba(255,255,255,0.04)', borderRadius: '0.25rem', width: '50%' }} />
                </div>
              </div>
            ))}
          </div>
        )}

        {!loading && error && (
          <div style={{ textAlign: 'center', padding: '4rem 1rem' }}>
            <p style={{ fontSize: '1.1rem', marginBottom: '0.5rem', color: 'rgba(255,255,255,0.7)' }}>Failed to load events</p>
            <p style={{ fontSize: '0.9rem', color: 'rgba(255,255,255,0.4)', marginBottom: '1.5rem' }}>{error}</p>
            <button onClick={fetchEvents} style={{ padding: '0.625rem 1.5rem', background: '#F59E0B', color: '#000', border: 'none', borderRadius: '0.5rem', fontWeight: 600, cursor: 'pointer' }}>Try Again</button>
          </div>
        )}

        {!loading && !error && events.length === 0 && (
          <div style={{ textAlign: 'center', padding: '5rem 1rem' }}>
            <h3 style={{ fontSize: '1.25rem', fontWeight: 600, marginBottom: '0.5rem', color: 'rgba(255,255,255,0.7)' }}>No Events Found</h3>
            <p style={{ color: 'rgba(255,255,255,0.4)' }}>Check back soon for upcoming events</p>
          </div>
        )}

        {!loading && !error && events.length > 0 && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
            {events.map(event => {
              const img = getEventImage(event);
              const date = getEventDate(event);
              const venue = getVenueInfo(event);
              const tiers = getTiers(event);
              return (
                <Link key={(event as any).id} to={`/events/${(event as any).id}`} style={{ textDecoration: 'none', display: 'block' }}>
                  <div
                    style={{
                      background: 'rgba(255,255,255,0.04)', borderRadius: '1rem', overflow: 'hidden',
                      border: '1px solid rgba(255,255,255,0.07)', transition: 'all 0.2s ease', cursor: 'pointer',
                    }}
                    onMouseEnter={e => {
                      const el = e.currentTarget as HTMLElement;
                      el.style.transform = 'translateY(-4px)';
                      el.style.borderColor = 'rgba(245,158,11,0.3)';
                      el.style.boxShadow = '0 12px 40px rgba(0,0,0,0.4)';
                    }}
                    onMouseLeave={e => {
                      const el = e.currentTarget as HTMLElement;
                      el.style.transform = 'translateY(0)';
                      el.style.borderColor = 'rgba(255,255,255,0.07)';
                      el.style.boxShadow = 'none';
                    }}
                  >
                    <div style={{ position: 'relative', height: '200px', overflow: 'hidden' }}>
                      {img ? (
                        <img src={img} alt={(event as any).name || (event as any).title || 'Event'} style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                          onError={e => { (e.target as HTMLImageElement).style.display = 'none'; }} />
                      ) : (
                        <div style={{ width: '100%', height: '100%', background: 'linear-gradient(135deg, #1a2744, #0f1729)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                          <svg style={{ width: '3rem', height: '3rem', color: 'rgba(245,158,11,0.3)' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
                          </svg>
                        </div>
                      )}
                      {(event as any).status && (
                        <div style={{
                          position: 'absolute', top: '0.75rem', right: '0.75rem',
                          padding: '0.25rem 0.625rem', borderRadius: '999px',
                          fontSize: '0.7rem', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em',
                          background: (event as any).status === 'on_sale' ? 'rgba(34,197,94,0.9)' :
                            (event as any).status === 'published' ? 'rgba(245,158,11,0.9)' :
                            (event as any).status === 'sold_out' ? 'rgba(239,68,68,0.9)' : 'rgba(0,0,0,0.6)',
                          color: '#fff',
                        }}>
                          {(event as any).status === 'on_sale' ? 'On Sale' :
                           (event as any).status === 'published' ? 'Available' :
                           (event as any).status === 'sold_out' ? 'Sold Out' : (event as any).status}
                        </div>
                      )}
                    </div>
                    <div style={{ padding: '1.25rem' }}>
                      <h3 style={{
                        fontFamily: "'Playfair Display', serif", fontSize: '1.05rem', fontWeight: 700,
                        color: '#fff', marginBottom: '0.5rem', lineHeight: 1.3,
                        display: '-webkit-box', WebkitLineClamp: 2, WebkitBoxOrient: 'vertical', overflow: 'hidden',
                      }}>{(event as any).name || (event as any).title}</h3>
                      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.35rem', marginBottom: '1rem' }}>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'rgba(255,255,255,0.5)', fontSize: '0.82rem' }}>
                          <svg style={{ width: '0.85rem', height: '0.85rem', flexShrink: 0 }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                          </svg>
                          {formatDate(date)}
                        </div>
                        <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: 'rgba(255,255,255,0.5)', fontSize: '0.82rem' }}>
                          <svg style={{ width: '0.85rem', height: '0.85rem', flexShrink: 0 }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                          </svg>
                          {venue.name}{venue.city ? `, ${venue.city}` : ''}
                        </div>
                      </div>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span style={{ color: '#F59E0B', fontWeight: 700, fontSize: '0.95rem' }}>{formatPrice(tiers)}</span>
                        <span style={{
                          padding: '0.35rem 0.875rem', background: 'rgba(245,158,11,0.1)',
                          border: '1px solid rgba(245,158,11,0.3)', borderRadius: '0.5rem',
                          color: '#F59E0B', fontSize: '0.8rem', fontWeight: 600,
                        }}>Get Tickets</span>
                      </div>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        )}

        {!loading && totalPages > 1 && (
          <div style={{ display: 'flex', justifyContent: 'center', gap: '0.5rem', marginTop: '3rem' }}>
            <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}
              style={{ padding: '0.5rem 1rem', background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', borderRadius: '0.5rem', color: page === 1 ? 'rgba(255,255,255,0.3)' : 'rgba(255,255,255,0.7)', cursor: page === 1 ? 'not-allowed' : 'pointer', fontSize: '0.875rem' }}>Previous</button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map(p => (
              <button key={p} onClick={() => setPage(p)}
                style={{ padding: '0.5rem 0.875rem', background: page === p ? '#F59E0B' : 'rgba(255,255,255,0.06)', border: '1px solid', borderColor: page === p ? '#F59E0B' : 'rgba(255,255,255,0.12)', borderRadius: '0.5rem', color: page === p ? '#000' : 'rgba(255,255,255,0.7)', cursor: 'pointer', fontWeight: page === p ? 700 : 400, fontSize: '0.875rem' }}>{p}</button>
            ))}
            <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages}
              style={{ padding: '0.5rem 1rem', background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.12)', borderRadius: '0.5rem', color: page === totalPages ? 'rgba(255,255,255,0.3)' : 'rgba(255,255,255,0.7)', cursor: page === totalPages ? 'not-allowed' : 'pointer', fontSize: '0.875rem' }}>Next</button>
          </div>
        )}
      </div>
    </div>
  );
}

/*
 * AdminEventCreatePage — uduXPass Admin
 * Design: Dark navy/amber brand system
 * Features: Hero image upload, gallery, promo video, per-tier ticket artwork
 */

import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Trash2, ArrowLeft, Calendar, MapPin, Ticket, CreditCard, Image as ImageIcon } from 'lucide-react'
import MediaUploader from '../../components/admin/MediaUploader'
import GalleryUploader from '../../components/admin/GalleryUploader'

interface Category {
  id: string; name: string; slug: string; description: string
}
interface Organizer { id: string; name: string; email: string; is_active: boolean }
interface TicketTier { name: string; description: string; price: string; quantity: string; maxPerOrder: number; imageUrl: string }

const AdminEventCreatePage: React.FC = () => {
  const navigate = useNavigate()
  const [categories, setCategories] = useState<Category[]>([])
  const [loadingCategories, setLoadingCategories] = useState(true)
  const [organizers, setOrganizers] = useState<Organizer[]>([])
  const [selectedOrganizerId, setSelectedOrganizerId] = useState('')
  const [loadingOrganizers, setLoadingOrganizers] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState('')

  const getDateTime = (daysFromNow: number, hour = 20) => {
    const d = new Date(); d.setDate(d.getDate() + daysFromNow); d.setHours(hour, 0, 0, 0)
    return d.toISOString().slice(0, 16)
  }

  const [eventData, setEventData] = useState({
    title: '', description: '', category: '',
    startDate: getDateTime(30), endDate: getDateTime(31, 0),
    venueName: '', venueAddress: '', venueCity: '', venueState: '', venueCountry: 'Nigeria', venueCapacity: '',
    eventImageUrl: '', galleryImages: [] as string[], promoVideoUrl: '',
    ticketTiers: [{ name: '', description: '', price: '', quantity: '', maxPerOrder: 10, imageUrl: '' }] as TicketTier[],
    enableMomo: true, enablePaystack: true
  })

  useEffect(() => {
    const token = localStorage.getItem('adminToken')
    fetch('/v1/categories').then(r => r.json()).then(d => {
      if (d.success && d.data) setCategories(d.data)
    }).catch(console.error).finally(() => setLoadingCategories(false))
    fetch('/v1/admin/organizers', { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json()).then(d => {
        if (d.success && d.data?.length > 0) { setOrganizers(d.data); setSelectedOrganizerId(d.data[0].id) }
      }).catch(console.error).finally(() => setLoadingOrganizers(false))
  }, [])

  const setField = (field: string, value: string | number | boolean) => {
    setEventData(prev => {
      const u = { ...prev, [field]: value }
      if (field === 'startDate') {
        try {
          const s = new Date(value as string); s.setDate(s.getDate() + 1); s.setHours(0, 0, 0, 0)
          u.endDate = s.toISOString().slice(0, 16)
        } catch { /* ignore */ }
      }
      return u
    })
  }

  const addTier = () => setEventData(p => ({ ...p, ticketTiers: [...p.ticketTiers, { name: '', description: '', price: '', quantity: '', maxPerOrder: 10, imageUrl: '' }] }))
  const removeTier = (i: number) => { if (eventData.ticketTiers.length > 1) setEventData(p => ({ ...p, ticketTiers: p.ticketTiers.filter((_, j) => j !== i) })) }
  const setTier = (i: number, f: string, v: string | number) => setEventData(p => ({ ...p, ticketTiers: p.ticketTiers.map((t, j) => j === i ? { ...t, [f]: v } : t) }))

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); setSubmitError('')
    if (!eventData.title.trim()) { setSubmitError('Please enter an event title'); return }
    if (!eventData.category) { setSubmitError('Please select a category'); return }
    if (!eventData.venueName.trim()) { setSubmitError('Please enter a venue name'); return }
    if (!eventData.venueCity.trim()) { setSubmitError('Please enter a city'); return }
    if (!selectedOrganizerId) { setSubmitError('No organizer available. Please create an organizer first.'); return }
    const validTiers = eventData.ticketTiers.filter(t => t.name.trim() && t.price && parseFloat(t.price) > 0)
    if (!validTiers.length) { setSubmitError('Please add at least one valid ticket tier with name and price'); return }
    setSubmitting(true)
    try {
      const token = localStorage.getItem('adminToken')
      const body: Record<string, unknown> = {
        organizer_id: selectedOrganizerId,
        name: eventData.title,
        slug: eventData.title.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '') + '-' + Date.now().toString(36),
        description: eventData.description || undefined,
        event_date: new Date(eventData.startDate).toISOString(),
        venue_name: eventData.venueName,
        venue_address: eventData.venueAddress || eventData.venueName,
        venue_city: eventData.venueCity,
        venue_state: eventData.venueState || undefined,
        venue_country: eventData.venueCountry || 'Nigeria',
        venue_capacity: parseInt(eventData.venueCapacity) || undefined,
        ticket_tiers: validTiers.map((t, i) => ({
          name: t.name.trim(), price: parseFloat(t.price), quota: parseInt(t.quantity) || 100,
          max_per_order: t.maxPerOrder || 10, description: t.description || undefined,
          sort_order: i + 1, image_url: t.imageUrl || undefined
        })),
        enable_momo: eventData.enableMomo, enable_paystack: eventData.enablePaystack
      }
      if (eventData.eventImageUrl) body.event_image_url = eventData.eventImageUrl
      if (eventData.promoVideoUrl) body.promo_video_url = eventData.promoVideoUrl
      if (eventData.galleryImages.length) body.gallery_images = eventData.galleryImages

      const res = await fetch('/v1/admin/events', {
        method: 'POST', headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      })
      if (res.ok) { navigate('/admin/events') }
      else { const err = await res.json().catch(() => ({})); setSubmitError(err.error || err.message || 'Failed to create event') }
    } catch { setSubmitError('An error occurred while creating the event.') }
    finally { setSubmitting(false) }
  }

  const iS: React.CSSProperties = { background: 'var(--bg-card)', borderColor: 'var(--border-color)', color: 'var(--text-primary)' }
  const iC = 'w-full rounded-lg px-3 py-2.5 text-sm border focus:outline-none transition-colors'
  const lC = 'block text-sm font-medium mb-1.5'
  const sC = 'rounded-2xl border p-6 space-y-5'

  return (
    <div className="min-h-screen" style={{ background: 'var(--bg-primary)', color: 'var(--text-primary)' }}>
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="mb-8">
          <button onClick={() => navigate('/admin/events')} className="flex items-center gap-2 text-sm mb-4 transition-colors hover:opacity-80" style={{ color: 'var(--text-secondary)' }}>
            <ArrowLeft size={16} /> Back to Events
          </button>
          <h1 className="text-3xl font-bold" style={{ fontFamily: 'var(--font-display)' }}>Create New Event</h1>
          <p className="mt-1 text-sm" style={{ color: 'var(--text-secondary)' }}>Set up your event with media, venue details, and ticket tiers</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">

          {/* Event Details */}
          <div className={sC} style={{ borderColor: 'var(--border-color)', background: 'var(--bg-elevated)' }}>
            <div className="flex items-center gap-2 mb-1">
              <Calendar size={18} style={{ color: 'var(--accent-amber)' }} />
              <h2 className="text-lg font-semibold">Event Details</h2>
            </div>
            <div>
              <label className={lC}>Event Title *</label>
              <input type="text" className={iC} style={iS} placeholder="e.g. Wizkid Live in Lagos" value={eventData.title} onChange={e => setField('title', e.target.value)} required />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className={lC}>Category *</label>
                {loadingCategories ? <div className="h-10 rounded-lg animate-pulse" style={{ background: 'var(--bg-card)' }} /> :
                  <select className={iC} style={iS} value={eventData.category} onChange={e => setField('category', e.target.value)} required>
                    <option value="">Select category</option>
                    {categories.map(c => <option key={c.id} value={c.name}>{c.name}</option>)}
                    {['Music', 'Sports', 'Comedy', 'Arts & Culture', 'Business', 'Food & Drink', 'Technology', 'Fashion'].map(c => (
                      <option key={c} value={c}>{c}</option>
                    ))}
                  </select>}
              </div>
              <div>
                <label className={lC}>Organizer *</label>
                {loadingOrganizers ? <div className="h-10 rounded-lg animate-pulse" style={{ background: 'var(--bg-card)' }} /> :
                  organizers.length > 0 ?
                    <select className={iC} style={iS} value={selectedOrganizerId} onChange={e => setSelectedOrganizerId(e.target.value)} required>
                      {organizers.map(o => <option key={o.id} value={o.id}>{o.name}</option>)}
                    </select> :
                    <div className="rounded-lg px-3 py-2.5 text-sm border" style={{ ...iS, color: '#ef4444' }}>No organizers found — create one first</div>}
              </div>
            </div>
            <div>
              <label className={lC}>Description</label>
              <textarea className={iC} style={iS} rows={4} placeholder="Describe the event, performers, experience..." value={eventData.description} onChange={e => setField('description', e.target.value)} />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className={lC}>Start Date &amp; Time *</label>
                <input type="datetime-local" className={iC} style={iS} value={eventData.startDate} onChange={e => setField('startDate', e.target.value)} required />
              </div>
              <div>
                <label className={lC}>End Date &amp; Time</label>
                <input type="datetime-local" className={iC} style={iS} value={eventData.endDate} onChange={e => setField('endDate', e.target.value)} />
              </div>
            </div>
          </div>

          {/* Media */}
          <div className={sC} style={{ borderColor: 'var(--border-color)', background: 'var(--bg-elevated)' }}>
            <div className="flex items-center gap-2 mb-1">
              <ImageIcon size={18} style={{ color: 'var(--accent-amber)' }} />
              <h2 className="text-lg font-semibold">Event Media</h2>
              <span className="text-xs px-2 py-0.5 rounded-full ml-1" style={{ background: 'rgba(245,158,11,0.15)', color: 'var(--accent-amber)' }}>Recommended</span>
            </div>
            <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
              High-quality images significantly increase ticket sales. Upload files directly or paste hosted URLs (Cloudinary, GCS, etc.).
            </p>
            <MediaUploader
              label="Hero / Banner Image"
              value={eventData.eventImageUrl}
              onChange={url => setField('eventImageUrl', url)}
              category="event_image"
              accept="image/jpeg,image/png,image/webp"
              maxSizeMB={10}
              placeholder="https://storage.googleapis.com/bucket/event-banner.jpg"
              hint="Recommended: 1920×1080px (16:9). Main image shown on event cards and the detail page hero."
            />
            <GalleryUploader
              images={eventData.galleryImages}
              onChange={imgs => setEventData(p => ({ ...p, galleryImages: imgs }))}
              maxImages={8}
            />
            <MediaUploader
              label="Promo Video"
              value={eventData.promoVideoUrl}
              onChange={url => setField('promoVideoUrl', url)}
              category="video"
              accept="video/mp4,video/webm"
              maxSizeMB={100}
              placeholder="https://youtube.com/watch?v=... or https://vimeo.com/..."
              hint="YouTube, Vimeo, or direct MP4/WebM URL. Displayed as an embedded player on the event page."
            />
          </div>

          {/* Venue */}
          <div className={sC} style={{ borderColor: 'var(--border-color)', background: 'var(--bg-elevated)' }}>
            <div className="flex items-center gap-2 mb-1">
              <MapPin size={18} style={{ color: 'var(--accent-amber)' }} />
              <h2 className="text-lg font-semibold">Venue Details</h2>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className={lC}>Venue Name *</label>
                <input type="text" className={iC} style={iS} placeholder="e.g. Eko Atlantic City" value={eventData.venueName} onChange={e => setField('venueName', e.target.value)} required />
              </div>
              <div>
                <label className={lC}>Capacity</label>
                <input type="number" className={iC} style={iS} placeholder="e.g. 50000" value={eventData.venueCapacity} onChange={e => setField('venueCapacity', e.target.value)} />
              </div>
            </div>
            <div>
              <label className={lC}>Street Address</label>
              <input type="text" className={iC} style={iS} placeholder="e.g. Eko Atlantic City, Victoria Island" value={eventData.venueAddress} onChange={e => setField('venueAddress', e.target.value)} />
            </div>
            <div className="grid grid-cols-3 gap-4">
              <div><label className={lC}>City *</label><input type="text" className={iC} style={iS} placeholder="Lagos" value={eventData.venueCity} onChange={e => setField('venueCity', e.target.value)} required /></div>
              <div><label className={lC}>State</label><input type="text" className={iC} style={iS} placeholder="Lagos" value={eventData.venueState} onChange={e => setField('venueState', e.target.value)} /></div>
              <div><label className={lC}>Country</label><input type="text" className={iC} style={iS} placeholder="Nigeria" value={eventData.venueCountry} onChange={e => setField('venueCountry', e.target.value)} /></div>
            </div>
          </div>

          {/* Ticket Tiers */}
          <div className={sC} style={{ borderColor: 'var(--border-color)', background: 'var(--bg-elevated)' }}>
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2">
                <Ticket size={18} style={{ color: 'var(--accent-amber)' }} />
                <h2 className="text-lg font-semibold">Ticket Tiers</h2>
              </div>
              <button type="button" onClick={addTier} className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium" style={{ background: 'rgba(245,158,11,0.15)', color: 'var(--accent-amber)' }}>
                <Plus size={14} /> Add Tier
              </button>
            </div>
            <div className="space-y-4">
              {eventData.ticketTiers.map((tier, i) => (
                <div key={i} className="rounded-xl border p-4 space-y-4" style={{ borderColor: 'var(--border-color)', background: 'var(--bg-card)' }}>
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-semibold" style={{ color: 'var(--accent-amber)' }}>Tier {i + 1}</span>
                    {eventData.ticketTiers.length > 1 && (
                      <button type="button" onClick={() => removeTier(i)} className="p-1.5 rounded-lg hover:opacity-80" style={{ color: '#ef4444' }}>
                        <Trash2 size={14} />
                      </button>
                    )}
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div><label className={lC}>Tier Name *</label><input type="text" className={iC} style={iS} placeholder="e.g. VIP Diamond" value={tier.name} onChange={e => setTier(i, 'name', e.target.value)} /></div>
                    <div><label className={lC}>Price (₦) *</label><input type="number" className={iC} style={iS} placeholder="e.g. 75000" value={tier.price} onChange={e => setTier(i, 'price', e.target.value)} /></div>
                  </div>
                  <div>
                    <label className={lC}>Description</label>
                    <textarea className={iC} style={iS} rows={2} placeholder="Describe what's included in this tier..." value={tier.description} onChange={e => setTier(i, 'description', e.target.value)} />
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div><label className={lC}>Quantity Available</label><input type="number" className={iC} style={iS} placeholder="100" value={tier.quantity} onChange={e => setTier(i, 'quantity', e.target.value)} /></div>
                    <div><label className={lC}>Max Per Order</label><input type="number" className={iC} style={iS} placeholder="10" value={tier.maxPerOrder} onChange={e => setTier(i, 'maxPerOrder', parseInt(e.target.value) || 10)} /></div>
                  </div>
                  <MediaUploader
                    label="Ticket Artwork (optional)"
                    value={tier.imageUrl}
                    onChange={url => setTier(i, 'imageUrl', url)}
                    category="ticket_image"
                    accept="image/jpeg,image/png,image/webp"
                    maxSizeMB={5}
                    placeholder="https://storage.googleapis.com/bucket/vip-ticket.jpg"
                    hint="Unique artwork for this tier. Shown during checkout and on the digital ticket."
                  />
                </div>
              ))}
            </div>
          </div>

          {/* Payment Methods */}
          <div className={sC} style={{ borderColor: 'var(--border-color)', background: 'var(--bg-elevated)' }}>
            <div className="flex items-center gap-2 mb-1">
              <CreditCard size={18} style={{ color: 'var(--accent-amber)' }} />
              <h2 className="text-lg font-semibold">Payment Methods</h2>
            </div>
            <div className="space-y-3">
              {([
                ['enableMomo', 'MoMo PSB Payment', 'Enable mobile money payments via MoMo PSB', eventData.enableMomo],
                ['enablePaystack', 'Paystack Payment', 'Enable card and bank payments via Paystack', eventData.enablePaystack]
              ] as [string, string, string, boolean][]).map(([field, label, desc, checked]) => (
                <div key={field} className="flex items-center justify-between p-4 rounded-xl border transition-colors"
                  style={{ borderColor: checked ? 'var(--accent-amber)' : 'var(--border-color)', background: checked ? 'rgba(245,158,11,0.05)' : 'var(--bg-card)' }}>
                  <div>
                    <p className="font-medium text-sm">{label}</p>
                    <p className="text-xs mt-0.5" style={{ color: 'var(--text-secondary)' }}>{desc}</p>
                  </div>
                  <label className="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" className="sr-only" checked={checked} onChange={e => setField(field, e.target.checked)} />
                    <div className="w-11 h-6 rounded-full relative transition-colors" style={{ background: checked ? 'var(--accent-amber)' : 'var(--border-color)' }}>
                      <div className="absolute top-0.5 left-0.5 bg-white rounded-full h-5 w-5 transition-transform" style={{ transform: checked ? 'translateX(20px)' : 'translateX(0)' }} />
                    </div>
                  </label>
                </div>
              ))}
            </div>
            {!eventData.enableMomo && !eventData.enablePaystack && (
              <div className="p-3 rounded-xl" style={{ background: 'rgba(234,179,8,0.1)', border: '1px solid rgba(234,179,8,0.3)' }}>
                <p className="text-sm" style={{ color: '#fbbf24' }}>⚠️ At least one payment method should be enabled.</p>
              </div>
            )}
          </div>

          {submitError && (
            <div className="p-4 rounded-xl" style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.3)' }}>
              <p className="text-sm" style={{ color: '#f87171' }}>{submitError}</p>
            </div>
          )}

          <div className="flex justify-end gap-3 pb-8">
            <button type="button" onClick={() => navigate('/admin/events')} className="px-6 py-2.5 rounded-xl text-sm font-medium border transition-opacity hover:opacity-80"
              style={{ borderColor: 'var(--border-color)', color: 'var(--text-secondary)', background: 'transparent' }}>
              Cancel
            </button>
            <button type="submit" disabled={submitting} className="px-8 py-2.5 rounded-xl text-sm font-bold disabled:opacity-60 transition-opacity hover:opacity-90"
              style={{ background: 'var(--accent-amber)', color: '#0f1729' }}>
              {submitting ? 'Creating...' : 'Create Event'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default AdminEventCreatePage

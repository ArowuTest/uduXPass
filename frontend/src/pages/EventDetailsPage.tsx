/*
 * EventDetailsPage — uduXPass User Frontend
 * Design: Dark navy (#0f1729) base, amber (#F59E0B) accents
 * FIXED: Uses event_image_url (not banner_image_url)
 * ADDED: Gallery grid, promo video embed, ticket tier artwork
 */

import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useToast } from '@/components/ui/use-toast';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import { useCart } from '../contexts/CartContext';
import { useAuth } from '../contexts/AuthContext';
import { eventsAPI } from '../services/api';
import { Event } from '../types/api';
import { formatCurrency, formatDate } from '../lib/utils';
import {
  Calendar, MapPin, ArrowLeft, Plus, Minus, ShoppingCart,
  Clock, Users, Share2, Heart, Ticket, Info, Play, Image as ImageIcon
} from 'lucide-react';
import { motion } from 'framer-motion';

interface TicketSelection {
  tierId: string;
  tierName: string;
  quantity: number;
  price: number;
}

// Extend Event type locally to include new media fields
interface EventWithMedia extends Event {
  gallery_images?: string[];
  promo_video_url?: string;
  thumbnail_url?: string;
}

const EventDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { addItem } = useCart();
  const { isAuthenticated } = useAuth();

  const [event, setEvent] = useState<EventWithMedia | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [ticketSelections, setTicketSelections] = useState<{ [tierId: string]: number }>({});
  const [isFavorite, setIsFavorite] = useState<boolean>(false);
  const [lightboxImg, setLightboxImg] = useState<string | null>(null);

  useEffect(() => {
    if (id) loadEvent();
  }, [id]);

  const loadEvent = async (): Promise<void> => {
    if (!id) return;
    try {
      setIsLoading(true);
      setError(null);
      const response = await eventsAPI.getById(id);
      if (response.success && response.data) {
        setEvent(response.data as EventWithMedia);
      } else {
        setError(response.error || 'Failed to load event');
      }
    } catch {
      setError('Failed to load event details');
    } finally {
      setIsLoading(false);
    }
  };

  const formatTime = (dateString: string): string =>
    new Date(dateString).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true });

  const handleTicketQuantityChange = (tierId: string, quantity: number, maxQuantity: number) => {
    const newQuantity = Math.max(0, Math.min(quantity, Math.min(maxQuantity, 10)));
    setTicketSelections(prev => ({ ...prev, [tierId]: newQuantity }));
  };

  const getTotalTickets = (): number => Object.values(ticketSelections).reduce((s, q) => s + q, 0);

  const getTotalPrice = (): number => {
    if (!event?.ticket_tiers) return 0;
    return event.ticket_tiers.reduce((total, tier) => total + (tier.price * (ticketSelections[tier.id] || 0)), 0);
  };

  const getAvailableQuantity = (tier: any): number => Math.max(0, (tier.quota || 0) - (tier.sold || 0));

  const handleAddToCart = () => {
    if (!event) return;
    if (!isAuthenticated) {
      toast({ title: 'Sign in required', description: 'Please log in to purchase tickets', variant: 'destructive' });
      navigate('/login', { state: { from: `/events/${id}` } });
      return;
    }
    const selections = event.ticket_tiers.filter(t => ticketSelections[t.id] > 0);
    if (selections.length === 0) {
      toast({ title: 'No tickets selected', description: 'Please select at least one ticket', variant: 'destructive' });
      return;
    }
    selections.forEach(tier => addItem(event.id, tier, ticketSelections[tier.id]));
    toast({ title: 'Added to cart', description: `${getTotalTickets()} ticket(s) added` });
    navigate('/checkout');
  };

  const handleShare = async () => {
    if (navigator.share) {
      try { await navigator.share({ title: event?.name, text: event?.description, url: window.location.href }); }
      catch { /* user cancelled */ }
    } else {
      navigator.clipboard.writeText(window.location.href);
      toast({ title: 'Link copied', description: 'Event link copied to clipboard' });
    }
  };

  // Resolve video embed URL (YouTube, Vimeo, or direct)
  const getVideoEmbed = (url: string): string | null => {
    if (!url) return null;
    const ytMatch = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&\s]+)/);
    if (ytMatch) return `https://www.youtube.com/embed/${ytMatch[1]}`;
    const vimeoMatch = url.match(/vimeo\.com\/(\d+)/);
    if (vimeoMatch) return `https://player.vimeo.com/video/${vimeoMatch[1]}`;
    return url; // direct mp4/webm
  };

  // ─── Loading ────────────────────────────────────────────────────────────────
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center" style={{ background: 'var(--bg-primary)' }}>
        <div className="text-center">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-sm" style={{ color: 'var(--text-secondary)' }}>Loading event details…</p>
        </div>
      </div>
    );
  }

  // ─── Error ──────────────────────────────────────────────────────────────────
  if (error || !event) {
    return (
      <div className="min-h-screen flex items-center justify-center" style={{ background: 'var(--bg-primary)' }}>
        <div className="max-w-md w-full mx-4 text-center p-10 rounded-2xl" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
          <Ticket size={48} style={{ color: 'var(--text-secondary)', margin: '0 auto 1rem' }} />
          <h1 className="text-2xl font-bold mb-2" style={{ color: 'var(--text-primary)' }}>Event Not Found</h1>
          <p className="mb-6 text-sm" style={{ color: 'var(--text-secondary)' }}>{error || "The event you're looking for doesn't exist."}</p>
          <button onClick={() => navigate('/events')} className="btn-primary w-full flex items-center justify-center gap-2">
            <ArrowLeft size={16} /> Back to Events
          </button>
        </div>
      </div>
    );
  }

  const totalAvailableTickets = event.ticket_tiers?.reduce((s, t) => s + getAvailableQuantity(t), 0) || 0;
  const heroImage = event.event_image_url || (event as any).banner_image_url;
  const gallery: string[] = event.gallery_images || [];
  const videoEmbed = event.promo_video_url ? getVideoEmbed(event.promo_video_url) : null;

  return (
    <div className="min-h-screen" style={{ background: 'var(--bg-primary)' }}>

      {/* ── Hero ── */}
      <div className="relative overflow-hidden" style={{ height: '480px', background: 'linear-gradient(135deg, #1e2d4a 0%, #0f1729 100%)' }}>
        {heroImage && (
          <img src={heroImage} alt={event.name} className="absolute inset-0 w-full h-full object-cover" style={{ opacity: 0.45 }} />
        )}
        <div className="absolute inset-0" style={{ background: 'linear-gradient(to top, rgba(15,23,41,0.95) 0%, rgba(15,23,41,0.3) 60%, transparent 100%)' }} />

        {/* Back button */}
        <button
          onClick={() => navigate('/events')}
          className="absolute top-6 left-6 flex items-center gap-2 px-3 py-2 rounded-xl text-sm font-medium transition-all hover:opacity-80"
          style={{ background: 'rgba(255,255,255,0.12)', color: '#f1f5f9', backdropFilter: 'blur(8px)' }}
        >
          <ArrowLeft size={16} /> Back
        </button>

        {/* Action buttons */}
        <div className="absolute top-6 right-6 flex gap-2">
          <button
            onClick={() => setIsFavorite(!isFavorite)}
            className="p-2.5 rounded-xl transition-all hover:opacity-80"
            style={{ background: 'rgba(255,255,255,0.12)', color: isFavorite ? '#f87171' : '#f1f5f9', backdropFilter: 'blur(8px)' }}
          >
            <Heart size={18} fill={isFavorite ? 'currentColor' : 'none'} />
          </button>
          <button
            onClick={handleShare}
            className="p-2.5 rounded-xl transition-all hover:opacity-80"
            style={{ background: 'rgba(255,255,255,0.12)', color: '#f1f5f9', backdropFilter: 'blur(8px)' }}
          >
            <Share2 size={18} />
          </button>
        </div>

        {/* Hero content */}
        <div className="absolute bottom-0 left-0 right-0 px-6 pb-8 max-w-5xl mx-auto">
          <div className="flex flex-wrap gap-2 mb-4">
            <span className="text-xs font-semibold px-3 py-1 rounded-full" style={{ background: event.status === 'published' ? 'rgba(34,197,94,0.2)' : 'rgba(245,158,11,0.2)', color: event.status === 'published' ? '#4ade80' : '#F59E0B' }}>
              {event.status === 'published' ? 'On Sale' : event.status?.toUpperCase()}
            </span>
            <span className="text-xs font-semibold px-3 py-1 rounded-full" style={{ background: 'rgba(255,255,255,0.12)', color: '#f1f5f9' }}>
              {totalAvailableTickets} tickets available
            </span>
          </div>
          <h1 className="text-4xl md:text-5xl font-bold mb-5 leading-tight" style={{ color: '#f1f5f9', fontFamily: "'Playfair Display', serif" }}>
            {event.name}
          </h1>
          <div className="flex flex-wrap gap-6" style={{ color: 'rgba(241,245,249,0.8)' }}>
            <div className="flex items-center gap-2">
              <Calendar size={16} style={{ color: '#F59E0B' }} />
              <div>
                <p className="text-sm font-semibold text-white">{formatDate(event.event_date)}</p>
                <p className="text-xs">{formatTime(event.event_date)}</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <MapPin size={16} style={{ color: '#F59E0B' }} />
              <div>
                <p className="text-sm font-semibold text-white">{event.venue_name || event.venue?.name}</p>
                <p className="text-xs">{event.venue_city || event.venue?.city}, {event.venue_state || event.venue?.state}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* ── Main Content ── */}
      <div className="max-w-5xl mx-auto px-4 py-10">
        <div className="grid gap-8 lg:grid-cols-3">

          {/* Left column */}
          <div className="lg:col-span-2 space-y-8">

            {/* About */}
            <motion.section initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.4 }}>
              <div className="rounded-2xl p-6" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                  <Info size={18} style={{ color: '#F59E0B' }} /> About This Event
                </h2>
                <p className="leading-relaxed whitespace-pre-line text-sm" style={{ color: 'var(--text-secondary)' }}>
                  {event.description || 'No description provided.'}
                </p>
              </div>
            </motion.section>

            {/* Gallery */}
            {gallery.length > 0 && (
              <motion.section initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.4, delay: 0.05 }}>
                <div className="rounded-2xl p-6" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                  <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                    <ImageIcon size={18} style={{ color: '#F59E0B' }} /> Gallery
                  </h2>
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                    {gallery.map((img, i) => (
                      <button key={i} onClick={() => setLightboxImg(img)} className="relative rounded-xl overflow-hidden aspect-video group">
                        <img src={img} alt={`Gallery ${i + 1}`} className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105" />
                        <div className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center" style={{ background: 'rgba(15,23,41,0.5)' }}>
                          <ImageIcon size={24} style={{ color: '#F59E0B' }} />
                        </div>
                      </button>
                    ))}
                  </div>
                </div>
              </motion.section>
            )}

            {/* Promo Video */}
            {videoEmbed && (
              <motion.section initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.4, delay: 0.1 }}>
                <div className="rounded-2xl p-6" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                  <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                    <Play size={18} style={{ color: '#F59E0B' }} /> Promo Video
                  </h2>
                  <div className="relative rounded-xl overflow-hidden" style={{ paddingBottom: '56.25%' }}>
                    <iframe
                      src={videoEmbed}
                      className="absolute inset-0 w-full h-full"
                      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                      allowFullScreen
                      title="Event promo video"
                    />
                  </div>
                </div>
              </motion.section>
            )}

            {/* Venue */}
            <motion.section initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.4, delay: 0.15 }}>
              <div className="rounded-2xl p-6" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                <h2 className="text-lg font-bold mb-4 flex items-center gap-2" style={{ color: 'var(--text-primary)' }}>
                  <MapPin size={18} style={{ color: '#F59E0B' }} /> Venue
                </h2>
                <p className="font-semibold" style={{ color: 'var(--text-primary)' }}>{event.venue_name || event.venue?.name}</p>
                <p className="text-sm mt-1" style={{ color: 'var(--text-secondary)' }}>{event.venue_address || event.venue?.address}</p>
                <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                  {event.venue_city || event.venue?.city}, {event.venue_state || event.venue?.state}, {event.venue_country || event.venue?.country}
                </p>
                {(event.venue_capacity || event.venue?.capacity) && (
                  <div className="flex items-center gap-2 mt-3 text-sm" style={{ color: 'var(--text-secondary)' }}>
                    <Users size={14} style={{ color: '#F59E0B' }} />
                    Capacity: {(event.venue_capacity || event.venue?.capacity || 0).toLocaleString()}
                  </div>
                )}
              </div>
            </motion.section>

            {/* Stats row */}
            <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.4, delay: 0.2 }}
              className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[
                { icon: Ticket, label: 'Available', value: totalAvailableTickets },
                { icon: Users, label: 'Sold', value: (event as any).tickets_sold || 0 },
                { icon: Calendar, label: 'Tiers', value: event.ticket_tiers?.length || 0 },
                { icon: Clock, label: 'Min Hold', value: '10' },
              ].map(({ icon: Icon, label, value }) => (
                <div key={label} className="rounded-xl p-4 text-center" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                  <Icon size={22} style={{ color: '#F59E0B', margin: '0 auto 8px' }} />
                  <p className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>{value}</p>
                  <p className="text-xs" style={{ color: 'var(--text-secondary)' }}>{label}</p>
                </div>
              ))}
            </motion.div>
          </div>

          {/* ── Ticket Sidebar ── */}
          <div className="lg:col-span-1">
            <motion.div initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.4 }} className="sticky top-8">
              <div className="rounded-2xl overflow-hidden" style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)' }}>
                <div className="px-5 py-4" style={{ borderBottom: '1px solid var(--border-color)' }}>
                  <h3 className="font-bold text-base" style={{ color: 'var(--text-primary)' }}>Select Tickets</h3>
                </div>
                <div className="p-5 space-y-4">
                  {event.ticket_tiers && event.ticket_tiers.length > 0 ? (
                    <>
                      {event.ticket_tiers.map((tier) => {
                        const available = getAvailableQuantity(tier);
                        const isAvailable = available > 0;
                        const tierImage = (tier as any).image_url;
                        return (
                          <div key={tier.id} className="rounded-xl p-4 transition-all"
                            style={{
                              background: isAvailable ? 'var(--bg-primary)' : 'rgba(255,255,255,0.03)',
                              border: `1px solid ${isAvailable ? 'var(--border-color)' : 'rgba(255,255,255,0.05)'}`,
                              opacity: isAvailable ? 1 : 0.5,
                            }}>
                            {/* Tier artwork */}
                            {tierImage && (
                              <img src={tierImage} alt={tier.name} className="w-full h-24 object-cover rounded-lg mb-3" />
                            )}
                            <div className="flex justify-between items-start mb-2">
                              <div className="flex-1 min-w-0">
                                <h4 className="font-semibold text-sm" style={{ color: 'var(--text-primary)' }}>{tier.name}</h4>
                                {tier.description && (
                                  <p className="text-xs mt-0.5 line-clamp-2" style={{ color: 'var(--text-secondary)' }}>{tier.description}</p>
                                )}
                                <p className="text-xs mt-1" style={{ color: 'var(--text-secondary)' }}>
                                  {available} of {tier.quota || tier.quantity || 0} available
                                </p>
                              </div>
                              <div className="text-right ml-3">
                                <p className="font-bold text-base" style={{ color: '#F59E0B' }}>{formatCurrency(tier.price)}</p>
                              </div>
                            </div>
                            {isAvailable ? (
                              <div className="flex items-center justify-between pt-3" style={{ borderTop: '1px solid var(--border-color)' }}>
                                <span className="text-xs" style={{ color: 'var(--text-secondary)' }}>Qty:</span>
                                <div className="flex items-center gap-2">
                                  <button
                                    onClick={() => handleTicketQuantityChange(tier.id, (ticketSelections[tier.id] || 0) - 1, available)}
                                    disabled={!ticketSelections[tier.id]}
                                    className="w-7 h-7 rounded-lg flex items-center justify-center transition-all disabled:opacity-30"
                                    style={{ background: 'var(--bg-elevated)', border: '1px solid var(--border-color)', color: 'var(--text-primary)' }}
                                  >
                                    <Minus size={12} />
                                  </button>
                                  <span className="w-8 text-center font-bold text-sm" style={{ color: 'var(--text-primary)' }}>
                                    {ticketSelections[tier.id] || 0}
                                  </span>
                                  <button
                                    onClick={() => handleTicketQuantityChange(tier.id, (ticketSelections[tier.id] || 0) + 1, available)}
                                    disabled={(ticketSelections[tier.id] || 0) >= Math.min(available, 10)}
                                    className="w-7 h-7 rounded-lg flex items-center justify-center transition-all disabled:opacity-30"
                                    style={{ background: 'rgba(245,158,11,0.15)', border: '1px solid rgba(245,158,11,0.3)', color: '#F59E0B' }}
                                  >
                                    <Plus size={12} />
                                  </button>
                                </div>
                              </div>
                            ) : (
                              <div className="mt-3 pt-3 text-center text-xs font-semibold" style={{ borderTop: '1px solid var(--border-color)', color: '#f87171' }}>
                                Sold Out
                              </div>
                            )}
                          </div>
                        );
                      })}

                      {/* Cart summary */}
                      {getTotalTickets() > 0 && (
                        <div className="pt-4 space-y-3" style={{ borderTop: '1px solid var(--border-color)' }}>
                          <div className="space-y-2 text-sm">
                            <div className="flex justify-between">
                              <span style={{ color: 'var(--text-secondary)' }}>Subtotal ({getTotalTickets()} tickets)</span>
                              <span className="font-semibold" style={{ color: 'var(--text-primary)' }}>{formatCurrency(getTotalPrice())}</span>
                            </div>
                            <div className="flex justify-between">
                              <span style={{ color: 'var(--text-secondary)' }}>Service Fee (5%)</span>
                              <span className="font-semibold" style={{ color: 'var(--text-primary)' }}>{formatCurrency(getTotalPrice() * 0.05)}</span>
                            </div>
                            <div className="flex justify-between text-base font-bold pt-2" style={{ borderTop: '1px solid var(--border-color)', color: 'var(--text-primary)' }}>
                              <span>Total</span>
                              <span style={{ color: '#F59E0B' }}>{formatCurrency(getTotalPrice() * 1.05)}</span>
                            </div>
                          </div>
                          <button
                            onClick={handleAddToCart}
                            className="btn-primary w-full flex items-center justify-center gap-2 py-3 text-sm font-semibold"
                          >
                            <ShoppingCart size={16} /> Proceed to Checkout
                          </button>
                          <p className="text-xs text-center" style={{ color: 'var(--text-secondary)' }}>
                            Tickets held for 10 minutes during checkout
                          </p>
                        </div>
                      )}
                    </>
                  ) : (
                    <div className="text-center py-10">
                      <Ticket size={40} style={{ color: 'var(--text-secondary)', margin: '0 auto 12px' }} />
                      <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>No tickets available at this time</p>
                    </div>
                  )}
                </div>
              </div>
            </motion.div>
          </div>
        </div>
      </div>

      {/* ── Lightbox ── */}
      {lightboxImg && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center p-4"
          style={{ background: 'rgba(0,0,0,0.9)' }}
          onClick={() => setLightboxImg(null)}
        >
          <img src={lightboxImg} alt="Gallery" className="max-w-full max-h-full rounded-xl object-contain" />
        </div>
      )}
    </div>
  );
};

export default EventDetailsPage;

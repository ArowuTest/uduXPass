/**
 * CreateSession.tsx — uduXPass Scanner PWA
 * Design: Enterprise Amber/Slate — "Professional Event Operations"
 * - Event selector with full event details preview
 * - Optional gate/location field
 * - Full API integration
 */

import { useState, useEffect } from 'react';
import { useLocation } from 'wouter';
import { scannerApi, Event } from '@/lib/api';
import { ArrowLeft, Loader2, MapPin, Calendar, Scan, ChevronDown, AlertCircle, Clock } from 'lucide-react';
import { toast } from 'sonner';

export default function CreateSession() {
  const [, setLocation] = useLocation();
  const [events, setEvents] = useState<Event[]>([]);
  const [selectedEventId, setSelectedEventId] = useState('');
  const [locationValue, setLocationValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingEvents, setIsLoadingEvents] = useState(true);
  const [showDropdown, setShowDropdown] = useState(false);

  useEffect(() => {
    scannerApi.getEvents()
      .then(setEvents)
      .catch(() => toast.error('Failed to load events'))
      .finally(() => setIsLoadingEvents(false));
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedEventId) { toast.error('Please select an event'); return; }
    setIsLoading(true);
    try {
      await scannerApi.createSession({
        event_id: selectedEventId,
        ...(locationValue.trim() && { location: locationValue.trim() }),
      });
      toast.success('Session started!');
      setLocation('/scan');
    } catch (err: any) {
      toast.error(err?.response?.data?.message || 'Failed to create session');
    } finally {
      setIsLoading(false);
    }
  };

  const selectedEvent = events.find((e) => e.id === selectedEventId);

  const formatEventDate = (iso: string) => {
    try {
      return new Date(iso).toLocaleDateString('en-NG', {
        weekday: 'short', day: 'numeric', month: 'short', year: 'numeric',
      });
    } catch { return iso; }
  };

  return (
    <div className="min-h-screen flex flex-col" style={{ background: 'oklch(0.13 0.025 245)' }}>

      {/* Header */}
      <header
        className="px-4 py-3 flex items-center gap-3 sticky top-0 z-20"
        style={{
          background: 'oklch(0.16 0.025 245)',
          borderBottom: '1px solid oklch(0.24 0.025 245)',
        }}
      >
        <button
          onClick={() => setLocation('/dashboard')}
          className="w-9 h-9 rounded-xl flex items-center justify-center flex-shrink-0"
          style={{
            background: 'oklch(0.22 0.025 245)',
            border: '1px solid oklch(0.28 0.025 245)',
            color: 'oklch(0.70 0.015 240)',
          }}
        >
          <ArrowLeft size={16} />
        </button>
        <div>
          <h1
            style={{
              fontFamily: 'Space Grotesk, sans-serif',
              fontSize: '15px',
              fontWeight: 700,
              color: 'oklch(0.97 0.005 240)',
              lineHeight: 1.2,
            }}
          >
            New Scanning Session
          </h1>
          <p style={{ fontSize: '12px', color: 'oklch(0.50 0.015 240)', lineHeight: 1.2 }}>
            Select an event to begin
          </p>
        </div>
      </header>

      {/* Form */}
      <main className="flex-1 px-4 py-5 max-w-lg mx-auto w-full">
        <form onSubmit={handleSubmit} className="space-y-5">

          {/* Section: Event */}
          <div>
            <label
              className="block mb-2"
              style={{
                fontSize: '11px',
                fontWeight: 700,
                letterSpacing: '0.1em',
                textTransform: 'uppercase',
                color: 'oklch(0.55 0.015 240)',
              }}
            >
              Event <span style={{ color: '#F43F5E' }}>*</span>
            </label>

            {isLoadingEvents ? (
              <div
                className="h-14 rounded-2xl flex items-center justify-center gap-2"
                style={{ background: 'oklch(0.18 0.025 245)', border: '1px solid oklch(0.26 0.025 245)' }}
              >
                <Loader2 size={16} className="animate-spin" style={{ color: 'oklch(0.50 0.015 240)' }} />
                <span style={{ fontSize: '13px', color: 'oklch(0.50 0.015 240)' }}>Loading events...</span>
              </div>
            ) : events.length === 0 ? (
              <div
                className="h-14 rounded-2xl flex items-center justify-center gap-2"
                style={{
                  background: 'rgba(244,63,94,0.06)',
                  border: '1px solid rgba(244,63,94,0.2)',
                }}
              >
                <AlertCircle size={16} style={{ color: '#F43F5E' }} />
                <span style={{ fontSize: '13px', color: '#F43F5E' }}>No events assigned to you</span>
              </div>
            ) : (
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setShowDropdown(!showDropdown)}
                  className="w-full h-14 px-4 rounded-2xl flex items-center justify-between text-left"
                  style={{
                    background: 'oklch(0.18 0.025 245)',
                    border: `1px solid ${selectedEventId ? 'rgba(245,158,11,0.4)' : 'oklch(0.26 0.025 245)'}`,
                    transition: 'border-color 150ms ease',
                  }}
                >
                  <span
                    style={{
                      fontSize: '14px',
                      fontWeight: selectedEventId ? 600 : 400,
                      color: selectedEventId ? 'oklch(0.97 0.005 240)' : 'oklch(0.45 0.015 240)',
                    }}
                  >
                    {selectedEvent?.name || 'Select an event...'}
                  </span>
                  <ChevronDown
                    size={16}
                    style={{
                      color: 'oklch(0.50 0.015 240)',
                      transform: showDropdown ? 'rotate(180deg)' : 'rotate(0deg)',
                      transition: 'transform 200ms ease',
                    }}
                  />
                </button>

                {/* Dropdown */}
                {showDropdown && (
                  <div
                    className="absolute top-full left-0 right-0 mt-2 rounded-2xl overflow-hidden z-30"
                    style={{
                      background: 'oklch(0.20 0.025 245)',
                      border: '1px solid oklch(0.28 0.025 245)',
                      boxShadow: '0 16px 48px rgba(0,0,0,0.5)',
                    }}
                  >
                    {events.map((event, idx) => (
                      <button
                        key={event.id}
                        type="button"
                        onClick={() => { setSelectedEventId(event.id); setShowDropdown(false); }}
                        className="w-full px-4 py-3 text-left flex items-center justify-between"
                        style={{
                          borderBottom: idx < events.length - 1 ? '1px solid oklch(0.26 0.025 245)' : 'none',
                          background: selectedEventId === event.id ? 'rgba(245,158,11,0.08)' : 'transparent',
                          transition: 'background 100ms ease',
                        }}
                      >
                        <div>
                          <p style={{ fontSize: '14px', fontWeight: 600, color: 'oklch(0.97 0.005 240)' }}>
                            {event.name}
                          </p>
                          {event.location && (
                            <p style={{ fontSize: '12px', color: 'oklch(0.50 0.015 240)', marginTop: '2px' }}>
                              {event.location}
                            </p>
                          )}
                        </div>
                        {selectedEventId === event.id && (
                          <div
                            className="w-5 h-5 rounded-full flex items-center justify-center flex-shrink-0"
                            style={{ background: '#F59E0B' }}
                          >
                            <svg width="10" height="8" viewBox="0 0 10 8" fill="none">
                              <path d="M1 4L3.5 6.5L9 1" stroke="oklch(0.13 0.025 245)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                            </svg>
                          </div>
                        )}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Selected Event Preview */}
          {selectedEvent && (
            <div
              className="rounded-2xl p-4 space-y-2"
              style={{
                background: 'rgba(245,158,11,0.06)',
                border: '1px solid rgba(245,158,11,0.2)',
              }}
            >
              <p
                style={{
                  fontFamily: 'Space Grotesk, sans-serif',
                  fontSize: '15px',
                  fontWeight: 700,
                  color: 'oklch(0.97 0.005 240)',
                }}
              >
                {selectedEvent.name}
              </p>
              {selectedEvent.location && (
                <div className="flex items-center gap-2" style={{ color: 'oklch(0.60 0.015 240)', fontSize: '13px' }}>
                  <MapPin size={13} style={{ color: '#F59E0B', flexShrink: 0 }} />
                  <span>{selectedEvent.location}</span>
                </div>
              )}
              {selectedEvent.start_time && (
                <div className="flex items-center gap-2" style={{ color: 'oklch(0.60 0.015 240)', fontSize: '13px' }}>
                  <Calendar size={13} style={{ color: '#F59E0B', flexShrink: 0 }} />
                  <span>{formatEventDate(selectedEvent.start_time)}</span>
                </div>
              )}
            </div>
          )}

          {/* Gate / Location (optional) */}
          <div>
            <label
              className="block mb-2"
              style={{
                fontSize: '11px',
                fontWeight: 700,
                letterSpacing: '0.1em',
                textTransform: 'uppercase',
                color: 'oklch(0.55 0.015 240)',
              }}
            >
              Gate / Entrance{' '}
              <span style={{ fontSize: '10px', fontWeight: 400, color: 'oklch(0.42 0.015 240)', textTransform: 'none', letterSpacing: 0 }}>
                (optional)
              </span>
            </label>
            <div className="relative">
              <MapPin
                size={15}
                className="absolute left-4 top-1/2 -translate-y-1/2"
                style={{ color: 'oklch(0.45 0.015 240)' }}
              />
              <input
                type="text"
                value={locationValue}
                onChange={(e) => setLocationValue(e.target.value)}
                placeholder="e.g. Main Gate, VIP Entrance"
                className="input-field"
                style={{ paddingLeft: '40px' }}
                disabled={isLoading}
              />
            </div>
          </div>

          {/* Info box */}
          <div
            className="flex items-start gap-3 p-3 rounded-xl"
            style={{
              background: 'oklch(0.18 0.025 245)',
              border: '1px solid oklch(0.26 0.025 245)',
            }}
          >
            <Clock size={14} style={{ color: 'oklch(0.50 0.015 240)', flexShrink: 0, marginTop: '1px' }} />
            <p style={{ fontSize: '12px', color: 'oklch(0.55 0.015 240)', lineHeight: 1.5 }}>
              A scanning session tracks all ticket validations for a specific event. Only one active session is allowed at a time.
            </p>
          </div>

          {/* Submit */}
          <button
            type="submit"
            className="btn-brand"
            disabled={isLoading || isLoadingEvents || !selectedEventId}
          >
            {isLoading ? (
              <>
                <Loader2 size={20} className="animate-spin" />
                <span>Starting Session...</span>
              </>
            ) : (
              <>
                <Scan size={20} />
                <span>Start Scanning Session</span>
              </>
            )}
          </button>
        </form>
      </main>

      {/* Backdrop for dropdown */}
      {showDropdown && (
        <div
          className="fixed inset-0 z-20"
          onClick={() => setShowDropdown(false)}
        />
      )}
    </div>
  );
}

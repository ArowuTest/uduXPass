-- Migration 021: Add full media support to events and ticket_tiers
-- Adds: promo_video_url, gallery_images (JSONB array of URLs), thumbnail_url to events
-- Adds: image_url to ticket_tiers for per-tier artwork
-- All fields are nullable to maintain backward compatibility.

-- ─── events table ─────────────────────────────────────────────────────────────

-- thumbnail_url: auto-generated or manually set smaller version of event_image_url
-- Used on event cards (400x250) to avoid loading the full hero image in list views.
ALTER TABLE events
    ADD COLUMN IF NOT EXISTS thumbnail_url VARCHAR(500),
    ADD COLUMN IF NOT EXISTS promo_video_url VARCHAR(500),
    ADD COLUMN IF NOT EXISTS gallery_images JSONB NOT NULL DEFAULT '[]'::jsonb;

COMMENT ON COLUMN events.event_image_url IS 'Hero/banner image URL (full resolution, 1200x630 recommended)';
COMMENT ON COLUMN events.thumbnail_url IS 'Card thumbnail URL (400x250 recommended). Falls back to event_image_url if null.';
COMMENT ON COLUMN events.promo_video_url IS 'Promotional video URL. Supports YouTube, Vimeo embed URLs or direct MP4 CDN links.';
COMMENT ON COLUMN events.gallery_images IS 'JSON array of up to 8 additional image URLs for the event gallery. Example: ["https://cdn.../img1.jpg", "https://cdn.../img2.jpg"]';

-- ─── ticket_tiers table ───────────────────────────────────────────────────────

-- image_url: per-tier artwork (e.g. VIP ticket has gold design, GA has standard design)
ALTER TABLE ticket_tiers
    ADD COLUMN IF NOT EXISTS image_url VARCHAR(500);

COMMENT ON COLUMN ticket_tiers.image_url IS 'Ticket tier artwork URL. Displayed on the checkout page and on the digital ticket PDF.';

-- ─── Update existing seed events with placeholder images ──────────────────────
-- These are real Unsplash concert images used for development/demo purposes.
-- In production, organisers will upload their own images via the admin portal.

UPDATE events SET
    event_image_url = 'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=1200&q=80',
    thumbnail_url   = 'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=400&q=80',
    gallery_images  = '[
        "https://images.unsplash.com/photo-1501386761578-eac5c94b800a?w=800&q=80",
        "https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800&q=80",
        "https://images.unsplash.com/photo-1429962714451-bb934ecdc4ec?w=800&q=80",
        "https://images.unsplash.com/photo-1516450360452-9312f5e86fc7?w=800&q=80"
    ]'::jsonb
WHERE slug = 'davido-timeless-lagos-2026';

UPDATE events SET
    event_image_url = 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=1200&q=80',
    thumbnail_url   = 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400&q=80',
    gallery_images  = '[
        "https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=800&q=80",
        "https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800&q=80",
        "https://images.unsplash.com/photo-1506157786151-b8491531f063?w=800&q=80"
    ]'::jsonb
WHERE slug = 'wizkid-live-lagos-2026-test';

UPDATE events SET
    event_image_url = 'https://images.unsplash.com/photo-1571266028243-d220c6a3d8c2?w=1200&q=80',
    thumbnail_url   = 'https://images.unsplash.com/photo-1571266028243-d220c6a3d8c2?w=400&q=80'
WHERE slug = 'diag-concert-1000';

UPDATE events SET
    event_image_url = 'https://images.unsplash.com/photo-1598387993441-a364f854cfds?w=1200&q=80',
    thumbnail_url   = 'https://images.unsplash.com/photo-1598387993441-a364f854cfds?w=400&q=80'
WHERE slug = 'e2e-concert-1772274100';

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add category_id to events table
ALTER TABLE events ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES categories(id) ON DELETE SET NULL;

-- Create index for category_id
CREATE INDEX IF NOT EXISTS idx_events_category_id ON events(category_id);

-- Insert default categories
INSERT INTO categories (name, slug, description, icon, color, display_order) VALUES
('Music', 'music', 'Concerts, festivals, and live music performances', '🎵', '#FF6B6B', 1),
('Sports', 'sports', 'Sporting events, matches, and tournaments', '⚽', '#4ECDC4', 2),
('Arts & Theater', 'arts-theater', 'Theater, dance, opera, and performing arts', '🎭', '#95E1D3', 3),
('Comedy', 'comedy', 'Stand-up comedy and comedy shows', '😂', '#FFE66D', 4),
('Conferences', 'conferences', 'Business conferences, seminars, and workshops', '💼', '#A8E6CF', 5),
('Festivals', 'festivals', 'Cultural festivals and celebrations', '🎉', '#FFD3B6', 6),
('Food & Drink', 'food-drink', 'Food festivals, wine tastings, and culinary events', '🍽️', '#FFAAA5', 7),
('Nightlife', 'nightlife', 'Clubs, parties, and nightlife events', '🌃', '#FF8B94', 8),
('Family', 'family', 'Family-friendly events and activities', '👨‍👩‍👧‍👦', '#A8DADC', 9),
('Other', 'other', 'Other events and activities', '📅', '#B8B8D1', 10)
ON CONFLICT (slug) DO NOTHING;

-- Add comment to category_id column
COMMENT ON COLUMN events.category_id IS 'Foreign key reference to categories table';

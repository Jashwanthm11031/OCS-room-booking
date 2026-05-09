CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    block_id UUID REFERENCES blocks(id) ON DELETE CASCADE,
    room_name VARCHAR(100) NOT NULL,
    capacity INTEGER NOT NULL,
    is_available BOOLEAN DEFAULT true,
    allowed_purposes TEXT[] DEFAULT '{OA,Interview,PPT}',
    notes TEXT
);

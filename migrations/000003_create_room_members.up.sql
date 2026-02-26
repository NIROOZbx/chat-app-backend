CREATE TABLE IF NOT EXISTS room_members (
    room_id INT REFERENCES room(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member', 
    is_banned BOOLEAN DEFAULT FALSE, 
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (room_id, user_id)
);
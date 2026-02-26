

CREATE TABLE IF NOT EXISTS room (
    id SERIAL PRIMARY KEY,
    name varchar(100) NOT NULL,
    max_members int ,
    description varchar(255) NOT NULL,
    topic varchar(50) NOT NULL,
    invite_code varchar(50) UNIQUE,
    is_private boolean DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
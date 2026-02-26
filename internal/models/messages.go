package models

import "time"

type Message struct {
    ID        int       `db:"id"`
    RoomID    int       `db:"room_id"`
    UserID    int       `db:"user_id"`
    Content   string    `db:"content"`
    CreatedAt time.Time `db:"created_at"`
    Username string `db:"user_name"`
}
package models

import "time"

type User struct {
	ID           int       `db:"id" `
	UserName     string    `db:"user_name"`
	ProfileImage string    `db:"profile_image"`
	IsOnline     bool      `db:"is_online" `
	LastSeen     time.Time `db:"last_seen" `
	CreatedAt    time.Time `db:"created_at"`
}

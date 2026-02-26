package models

import "time"

type User struct {
	ID           int       `db:"id" json:"ID"`
	UserName     string    `db:"user_name" json:"UserName"`
	ProfileImage string    `db:"profile_image" json:"ProfileImage"`
	IsOnline     bool      `db:"is_online" json:"IsOnline"`
	LastSeen     time.Time `db:"last_seen" json:"LastSeen"`
	CreatedAt    time.Time `db:"created_at" json:"CreatedAt"`
}

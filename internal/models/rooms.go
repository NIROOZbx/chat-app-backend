package models

import (
	"time"
)

type Room struct {
	ID          int       `db:"id"`
	Name        string    `db:"name" `
	Description string    `db:"description"`
	Topic       string    `db:"topic"`
	IsPrivate   bool      `db:"is_private"`
	InviteCode  *string   `db:"invite_code"`
	Image       string    `db:"image"`

	MemberCount int64       `db:"-"`
	MaxMembers  int       `db:"max_members"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

package gonext

import "time"

// User represents a user
type User struct {
	ID           uint64     `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	LastSeen  time.Time `json:"last_seen" db:"last_seen"`
}

package session

import (
	"time"

	gonext "github.com/jackcaldwell/go-next"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to a session store
type Repository interface {
	Create(userID uint64) (*gonext.Session, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository returns a new instance of a session repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(userID uint64) (*gonext.Session, error) {
	s := gonext.Session{}

	tx := r.db.MustBegin()

	err := tx.QueryRowx(
		"INSERT INTO sessions (user_id, created_at, last_seen) VALUES ($1, $2, $3) RETURNING *",
		userID,
		time.Now(),
		time.Now(),
	).StructScan(&s)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	return &s, err
}

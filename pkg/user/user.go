package user

import (
	"time"

	gonext "github.com/jackcaldwell/go-next"
	"github.com/jackcaldwell/go-next/pkg/crypto"
	"github.com/jmoiron/sqlx"
)

// Repository provides access to a user store
type Repository interface {
	Create(email string, password string) (*gonext.User, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository returns a new instance of a user repository.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(email string, password string) (u *gonext.User, err error) {
	u = &gonext.User{}

	tx := r.db.MustBegin()

	err = tx.QueryRowx(
		"INSERT INTO users (email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING *",
		email,
		crypto.HashString(password),
		time.Now(),
		time.Now(),
	).StructScan(u)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	return u, err
}

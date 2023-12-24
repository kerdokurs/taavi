package data

import (
	"database/sql"
	"time"
)

type Model struct {
	ID uint `db:"id" json:"id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

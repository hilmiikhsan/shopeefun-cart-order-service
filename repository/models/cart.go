package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	ProductID uuid.UUID  `db:"product_id"`
	Qty       int        `db:"qty"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type GetCartRequest struct {
	UserID    uuid.UUID   `db:"user_id"`
	ProductID []uuid.UUID `db:"product_id"`
}

type CartRequest struct {
	UserID    uuid.UUID    `db:"user_id"`
	ProductID uuid.UUID    `db:"product_ids"`
	Qty       int          `db:"qty"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type DeleteCartRequest struct {
	UserID    uuid.UUID `db:"user_id"`
	ProductID uuid.UUID `db:"product_id"`
}

package models

import (
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

type AddCartRequest struct {
	UserID     uuid.UUID   `db:"user_id"`
	ProductIDs []uuid.UUID `db:"product_ids"`
	Qty        int         `db:"qty"`
}

type DeleteCartRequest struct {
	UserID    uuid.UUID   `db:"user_id"`
	ProductID []uuid.UUID `db:"product_id"`
}

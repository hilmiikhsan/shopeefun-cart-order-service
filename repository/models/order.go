package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UserID        uuid.UUID       `db:"user_id"`
	OrderID       uuid.UUID       `db:"order_id"`
	PaymentTypeID uuid.UUID       `db:"payment_type_id"`
	OrderNumber   string          `db:"order_number"`
	TotalPrice    float64         `db:"total_price"`
	ProductOrder  json.RawMessage `db:"product_order"`
	Status        string          `db:"status"`
	IsPaid        bool            `db:"is_paid"`
	RefCode       string          `db:"ref_code"`
	CreatedAt     *time.Time      `db:"created_at"`
	UpdatedAt     *time.Time      `db:"updated_at"`
	DeletedAt     *time.Time      `db:"deleted_at"`
}

type OrderItemsLogs struct {
	OrderID    uuid.UUID  `db:"order_id"`
	RefCode    string     `db:"ref_code"`
	FromStatus string     `db:"from_status"`
	ToStatus   string     `db:"to_status"`
	Notes      string     `db:"notes"`
	CreatedAt  *time.Time `db:"created_at"`
}

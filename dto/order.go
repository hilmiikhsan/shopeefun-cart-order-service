package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserID        uuid.UUID       `json:"user_id" validate:"required,uuid"`
	PaymentTypeID uuid.UUID       `json:"payment_type_id" validate:"required,uuid"`
	OrderNumber   string          `json:"order_number" validate:"required"`
	TotalPrice    float64         `json:"total_price" validate:"required"`
	ProductOrder  json.RawMessage `json:"product_order"`
	Status        string          `json:"status" validate:"required"`
	IsPaid        bool            `json:"is_paid"`
	RefCode       string          `json:"ref_code"`
	CreatedAt     *time.Time      `json:"created_at"`
	UpdatedAt     *time.Time      `json:"updated_at"`
	DeletedAt     *time.Time      `json:"deleted_at"`
}

type UpdateOrderRequest struct {
	OrderID   uuid.UUID  `json:"order_id" validate:"required,uuid"`
	Status    string     `json:"status" validate:"required"`
	IsPaid    bool       `json:"is_paid"`
	UpdatedAt *time.Time `json:"updated_at"`
}

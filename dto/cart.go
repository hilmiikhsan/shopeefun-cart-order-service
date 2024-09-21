package dto

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type GetCartRequest struct {
	UserID    string   `json:"user_id" validate:"required,uuid"`
	ProductID []string `json:"product_id" validate:"required,dive,uuid"`
}

type GetCartResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	ProductID string     `json:"product_id"`
	Qty       int        `json:"qty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type AddCartRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	ProductID string `json:"product_id" validate:"required,uuid"`
	Qty       int    `json:"qty" validate:"numeric"`
}

type ProductRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Qty       int    `json:"qty" validate:"numeric"`
}

type UpdateCartRequest struct {
	UserID   string           `json:"user_id" validate:"required,uuid"`
	Products []ProductRequest `json:"products" validate:"required,min=1,dive"`
}

type DeleteCartRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	ProductID string `json:"product_id" validate:"required,uuid"`
}

// Custom UnmarshalJSON method for ProductRequest
func (p *AddCartRequest) UnmarshalJSON(data []byte) error {
	type Alias AddCartRequest
	aux := &struct {
		// ProductID json.RawMessage `json:"product_id"`
		Qty json.RawMessage `json:"qty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert Qty from string or number to int
	var qty int
	if err := json.Unmarshal(aux.Qty, &qty); err != nil {
		var qtyStr string
		if err := json.Unmarshal(aux.Qty, &qtyStr); err != nil {
			return errors.New("qty harus angka")
		}
		var err error
		qty, err = strconv.Atoi(qtyStr)
		if err != nil {
			return errors.New("qty harus angka")
		}
	}
	p.Qty = qty

	return nil
}

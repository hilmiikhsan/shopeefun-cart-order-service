package dto

import (
	"encoding/json"
	"errors"
	"fmt"
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
	UserID     string   `json:"user_id" validate:"required,uuid"`
	ProductIDs []string `json:"product_ids" validate:"required,min=1,dive,uuid"`
	Qty        int      `json:"qty" validate:"required,numeric,min=1"`
}

type UpdateCartRequest struct {
	UserID   string           `json:"user_id" validate:"required,uuid"`
	Products []ProductRequest `json:"products" validate:"required,min=1,dive"`
}

type ProductRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Qty       int    `json:"qty" validate:"numeric"`
}

type DeleteCartRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	ProductID string `json:"product_id" validate:"required,uuid"`
}

// Custom UnmarshalJSON method for ProductRequest
func (p *ProductRequest) UnmarshalJSON(data []byte) error {
	type Alias ProductRequest
	aux := &struct {
		ProductID json.RawMessage `json:"product_id"`
		Qty       json.RawMessage `json:"qty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert ProductID from string or number to string
	var productIDStr string
	if err := json.Unmarshal(aux.ProductID, &productIDStr); err != nil {
		// If it's not a string, try converting from a number
		var productIDNum float64
		if err := json.Unmarshal(aux.ProductID, &productIDNum); err != nil {
			return errors.New("product_id must be a valid UUID or numeric ID")
		}
		productIDStr = fmt.Sprintf("%.0f", productIDNum)
	}
	p.ProductID = productIDStr

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

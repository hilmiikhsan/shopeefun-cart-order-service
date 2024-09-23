package order

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/repository/models"
	"github.com/sirupsen/logrus"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{db}
}

// CreateOrder is a method that creates a new order and returns the order ID.
// It returns an error if any occurs during the creation process.
func (o *store) CreateOrder(ctx context.Context, tx *sql.Tx, req models.Order) (*uuid.UUID, *string, error) {
	var orderID uuid.UUID
	var refCode string

	if err := tx.QueryRowContext(ctx, queryCreateOrder,
		req.UserID,
		req.PaymentTypeID,
		req.OrderNumber,
		req.TotalPrice,
		req.ProductOrder,
		req.Status,
		req.IsPaid,
		req.RefCode,
	).Scan(&orderID, &refCode); err != nil {
		logrus.Errorf("[Store][CreateOrder] Error on inserting new order: %v", err)
		return nil, nil, err
	}

	return &orderID, &refCode, nil
}

// CreateOrderItemsLogs is a method that creates a new order status logs.
// It returns an error if any occurs during the creation process.
func (o *store) CreateOrderItemsLogs(ctx context.Context, tx *sql.Tx, req models.OrderItemsLogs) (*string, error) {
	var refCode string

	if err := tx.QueryRowContext(ctx, queryCreateOrderStatusLogs,
		req.OrderID,
		req.RefCode,
		req.FromStatus,
		req.ToStatus,
		req.Notes,
	).Scan(&refCode); err != nil {
		logrus.Errorf("[Store][CreateOrderItemsLogs] Error on inserting new order status logs: %v", err)
		return nil, err
	}

	return &refCode, nil
}

// UpdateOrder is a method that updates an existing order.
// It returns an error if any occurs during the update process.
func (o *store) UpdateOrder(ctx context.Context, tx *sql.Tx, req models.Order) (*string, error) {
	var refCode string

	if err := tx.QueryRowContext(ctx, queryUpdateOrder,
		req.Status,
		req.IsPaid,
		req.OrderID,
	).Scan(&refCode); err != nil {
		logrus.Errorf("[Store][UpdateOrder] Error on updating order: %v", err)
		return nil, err
	}

	return &refCode, nil
}

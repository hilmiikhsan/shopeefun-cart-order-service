package order

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/constants"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/dto"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/repository/models"
	"github.com/sirupsen/logrus"
)

type orderStore interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, req models.Order) (*uuid.UUID, *string, error)
	CreateOrderItemsLogs(ctx context.Context, tx *sql.Tx, req models.OrderItemsLogs) (*string, error)
	UpdateOrder(ctx context.Context, tx *sql.Tx, req models.Order) (*string, error)
}

type order struct {
	ctx   context.Context
	store orderStore
	db    *sql.DB
}

func NewOrder(ctx context.Context, store orderStore, db *sql.DB) *order {
	return &order{
		ctx,
		store,
		db,
	}
}

func (o *order) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*uuid.UUID, error) {
	tx, err := o.db.Begin()
	if err != nil {
		logrus.Errorf("[Usecase][CreateOrder] Failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("[Usecase][CreateOrder] Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	orderID, refCode, err := o.store.CreateOrder(ctx, tx, models.Order{
		UserID:        req.UserID,
		PaymentTypeID: req.PaymentTypeID,
		OrderNumber:   req.OrderNumber,
		TotalPrice:    req.TotalPrice,
		ProductOrder:  req.ProductOrder,
		Status:        strings.ToLower(req.Status),
		IsPaid:        req.IsPaid,
		RefCode:       req.RefCode,
	})
	if err != nil {
		logrus.Errorf("[Usecase][CreateOrder] Failed to create order: %v", err)
		return nil, err
	}

	_, err = o.store.CreateOrderItemsLogs(ctx, tx, models.OrderItemsLogs{
		OrderID:    *orderID,
		RefCode:    *refCode,
		FromStatus: strings.ToLower(constants.OrderStatusPending),
		ToStatus:   strings.ToLower(constants.OrderStatusPaid),
		Notes:      constants.OrderCreatedsMessage,
	})
	if err != nil {
		logrus.Errorf("[Usecase][CreateOrder] Failed to create order items logs: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("[Usecase][CreateOrder] Failed to commit transaction: %v", err)
		return nil, err
	}

	return orderID, nil
}

func (o *order) UpdateOrder(ctx context.Context, req dto.UpdateOrderRequest) (*string, error) {
	tx, err := o.db.Begin()
	if err != nil {
		logrus.Errorf("[Usecase][UpdateOrder] Failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("[Usecase][UpdateOrder] Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	refCode, err := o.store.UpdateOrder(ctx, tx, models.Order{
		Status:  strings.ToLower(req.Status),
		IsPaid:  req.IsPaid,
		OrderID: req.OrderID,
	})
	if err != nil {
		logrus.Errorf("[Usecase][UpdateOrder] Failed to update order: %v", err)
		return nil, err
	}

	_, err = o.store.CreateOrderItemsLogs(ctx, tx, models.OrderItemsLogs{
		OrderID:    req.OrderID,
		RefCode:    *refCode,
		FromStatus: strings.ToLower(constants.OrderStatusPending),
		ToStatus:   strings.ToLower(constants.OrderStatusPaid),
		Notes:      strings.ToLower(constants.PaymentSuccessMessage),
	})
	if err != nil {
		logrus.Errorf("[Usecase][UpdateOrder] Failed to create order items logs: %v", err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("[Usecase][CreateOrder] Failed to commit transaction: %v", err)
		return nil, err
	}

	messageStatus := "Payment Success"

	return &messageStatus, nil
}

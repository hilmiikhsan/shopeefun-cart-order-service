package cart

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/repository/models"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/validators/errmsg"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type store struct {
	db *sql.DB
}

// NewStore is a constructor function that returns a new store instance.
func NewStore(db *sql.DB) *store {
	return &store{db}
}

// GetCartByUserID is a method that retrieves the cart for a given user.
// It returns a slice of cart and an error if any occurs during the retrieval process.
func (s *store) GetCartByUserID(ctx context.Context, req models.GetCartRequest) (*[]models.Cart, error) {
	query := queryGetAllCart
	var queryConditions []string

	if req.UserID != uuid.Nil {
		queryConditions = append(queryConditions, fmt.Sprintf("user_id = '%s'", req.UserID))
	}

	if len(req.ProductID) > 0 {
		var productIDs []string
		for _, pid := range req.ProductID {
			productIDs = append(productIDs, fmt.Sprintf("'%s'", pid))
		}
		queryConditions = append(queryConditions, fmt.Sprintf("product_id IN (%s)", strings.Join(productIDs, ",")))
	}

	if len(queryConditions) > 0 {
		query += " WHERE " + strings.Join(queryConditions, " AND ")
	} else {
		query += " WHERE deleted_at IS NULL"
	}

	query += " AND deleted_at IS NULL"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		logrus.Errorf("[Store][GetCartByUserID] Error on querying the database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var carts []models.Cart
	for rows.Next() {
		var cart models.Cart
		if err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.ProductID,
			&cart.Qty,
			&cart.CreatedAt,
			&cart.UpdatedAt,
			&cart.DeletedAt,
		); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}

	if err := rows.Err(); err != nil {
		logrus.Errorf("[Store][GetCartByUserID] error rows scan sql: %v", err)
		return nil, err
	}

	return &carts, nil
}

// AddCart is a method that adds new products to a user's cart.
// It returns the ID of the first inserted product and an error if any occurs during the addition process.
func (s *store) AddCart(ctx context.Context, req models.AddCartRequest) (*uuid.UUID, error) {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.Errorf("[Store][AddCart] Failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("[Store][AddCart] Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	var firstID uuid.UUID
	productIDStrings := make([]string, len(req.ProductIDs))
	for i, id := range req.ProductIDs {
		productIDStrings[i] = id.String()
	}

	productIDSet := make(map[string]bool)

	rows, err := tx.QueryContext(ctx, queryGetProductByUserIdAndProductId,
		req.UserID,
		pq.Array(productIDStrings),
	)
	if err != nil {
		logrus.Errorf("[Store][AddCart] Failed to query existing products: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			logrus.Errorf("[Store][AddCart] Failed to scan existing product: %v", err)
			return nil, err
		}
		productIDSet[id] = true
	}

	for _, productID := range req.ProductIDs {
		productIDStr := productID.String()
		if productIDSet[productIDStr] {
			continue
		}

		var id uuid.UUID
		if err := tx.QueryRowContext(ctx, queryAddCart,
			req.UserID,
			productIDStr,
			req.Qty,
		).Scan(&id); err != nil {
			logrus.Errorf("[Store][AddCart] Failed to insert cart: %v", err)
			return nil, err
		}

		if firstID == uuid.Nil {
			firstID = id
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("[Store][AddCart] Failed to commit transaction: %v", err)
		return nil, err
	}

	if firstID == uuid.Nil {
		return nil, fmt.Errorf(errmsg.ErrDoesntNewItemsAdded)
	}

	return &firstID, nil
}

// UpdateQty is a method that updates the quantity of a product in a user's cart.
// It returns an error if any occurs during the update process.
func (s *store) UpdateQty(ctx context.Context, userID, productID uuid.UUID, qty int) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("[Store][UpdateQty] Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	// Check if the user exists
	var userExists bool
	if err := tx.QueryRowContext(ctx, queryCheckUserExists, userID).Scan(&userExists); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to check user existence: %v", err)
		return err
	}
	if !userExists {
		return errors.New("user not found")
	}

	// Check if the product exists in the user's cart
	var productInCart bool
	if err := tx.QueryRowContext(ctx, queryCheckProductInCart,
		userID,
		productID,
	).Scan(&productInCart); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to check product in cart: %v", err)
		return err
	}
	if !productInCart {
		return errors.New("product not found in cart")
	}

	// Lock the cart for update
	if _, err := tx.ExecContext(ctx, queryLockUpdateQty, userID); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to lock cart: %v", err)
		return errors.New("failed to lock data")
	}

	// Update the product quantity
	if _, err := tx.ExecContext(ctx, queryUpdateQty,
		qty,
		userID,
		productID,
	); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to update cart: %v", err)
		return errors.New("failed to update data")
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

// DeleteProduct is a method that deletes a product from a user's cart.
// It returns an error if any occurs during the deletion process.
func (s *store) DeleteProduct(ctx context.Context, req models.DeleteCartRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		logrus.Errorf("[Store][DeleteProduct] Failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logrus.Errorf("[Store][UpdateQty] Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	// Check if the user exists
	var userExists bool
	if err := tx.QueryRowContext(ctx, queryCheckUserExists, req.UserID).Scan(&userExists); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to check user existence: %v", err)
		return err
	}
	if !userExists {
		return errors.New("user not found")
	}

	// Check if the product exists in the user's cart
	var productInCart bool
	if err := tx.QueryRowContext(ctx, queryCheckProductInCart,
		req.UserID,
		pq.Array(req.ProductID),
	).Scan(&productInCart); err != nil {
		logrus.Errorf("[Store][UpdateQty] Failed to check product in cart: %v", err)
		return err
	}
	if !productInCart {
		return errors.New("product not found in cart")
	}

	if _, err := tx.ExecContext(ctx, queryLockSoftDeleteProduct, req.UserID); err != nil {
		logrus.Errorf("[Store][DeleteProduct] Failed to lock cart: %v", err)
		return errors.New("failed to lock data")
	}

	if _, err := tx.ExecContext(ctx, queryUpdateDeletedAt, req.UserID, pq.Array(req.ProductID)); err != nil {
		logrus.Errorf("[Store][DeleteProduct] Failed to delete cart: %v", err)
		return errors.New("failed to delete data")
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorf("[Store][DeleteProduct] Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

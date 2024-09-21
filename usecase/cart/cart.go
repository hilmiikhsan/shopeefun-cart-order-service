package cart

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/dto"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/repository/models"
	"github.com/sirupsen/logrus"
)

// cartStore is an interface that defines the methods required for managing a shopping cart.
type cartStore interface {
	GetCartByUserID(ctx context.Context, req models.GetCartRequest) (*[]models.Cart, error)
	AddCart(ctx context.Context, req models.CartRequest) (models.Cart, error)
	UpdateQty(ctx context.Context, req models.CartRequest) (models.Cart, error)
	DeleteProduct(ctx context.Context, req models.DeleteCartRequest) error
}

// cart is a struct that holds the store for managing a shopping cart.
type cart struct {
	ctx   context.Context
	store cartStore
}

// NewCart is a constructor function that returns a new cart instance.
func NewCart(ctx context.Context, store cartStore) *cart {
	return &cart{
		ctx,
		store,
	}
}

// GetCartByUserID is a method that retrieves the cart for a given user and returns a response with the total items.
func (c *cart) GetCartByUserID(ctx context.Context, req dto.GetCartRequest) ([]dto.GetCartResponse, error) {
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		logrus.Errorf("[Handler][GetCartByUserID] Failed to parse user id: %v", err)
		return nil, err
	}

	var pidSlice []uuid.UUID
	for _, pid := range req.ProductID {
		parsedPID, err := uuid.Parse(pid)
		if err != nil {
			logrus.Errorf("[Handler][GetCartByUserID] Failed to parse product id: %v", err)
			return nil, err
		}
		pidSlice = append(pidSlice, parsedPID)
	}

	result, err := c.store.GetCartByUserID(ctx, models.GetCartRequest{
		UserID:    uid,
		ProductID: pidSlice,
	})
	if err != nil {
		log.Printf("[Usecase][GetCartByUserID] Error on getting cart by user id: %v", err)
		return nil, err
	}

	if len(*result) == 0 {
		return []dto.GetCartResponse{}, nil
	}

	cartResponse := make([]dto.GetCartResponse, 0)
	for _, cart := range *result {
		cartResponse = append(cartResponse, dto.GetCartResponse{
			ID:        cart.ID.String(),
			UserID:    cart.UserID.String(),
			ProductID: cart.ProductID.String(),
			Qty:       cart.Qty,
			CreatedAt: cart.CreatedAt,
			UpdatedAt: cart.UpdatedAt,
			DeletedAt: cart.DeletedAt,
		})
	}

	return cartResponse, nil
}

// AddCart is a method that adds a new product to a user's cart.
func (c *cart) AddCart(ctx context.Context, req dto.AddCartRequest) (dto.GetCartResponse, error) {
	parsedProductID, err := uuid.Parse(req.ProductID)
	if err != nil {
		logrus.Errorf("[Usecase][AddCart] Failed to parse product id: %v", err)
		return dto.GetCartResponse{}, fmt.Errorf("failed to parse product id %v: %w", req.ProductID, err)
	}

	cartData, err := c.store.AddCart(ctx, models.CartRequest{
		UserID:    uuid.MustParse(req.UserID),
		ProductID: parsedProductID,
		Qty:       req.Qty,
	})
	if err != nil {
		logrus.Errorf("[Usecase][AddCart] Error on adding cart: %v", err)
		return dto.GetCartResponse{}, err
	}

	return dto.GetCartResponse{
		ID:        cartData.ID.String(),
		UserID:    cartData.UserID.String(),
		ProductID: cartData.ProductID.String(),
		Qty:       cartData.Qty,
		CreatedAt: cartData.CreatedAt,
		UpdatedAt: cartData.UpdatedAt,
	}, nil
}

// UpdateQty is a method that updates the quantity of a product in a user's cart or deletes the product if the quantity is 0.
func (c *cart) UpdateQty(ctx context.Context, req dto.UpdateCartRequest) ([]dto.GetCartResponse, error) {
	var cartResponses []dto.GetCartResponse

	for _, product := range req.Products {
		if product.Qty == 0 {
			// Delete the product from the cart if the quantity is 0
			if err := c.store.DeleteProduct(ctx, models.DeleteCartRequest{
				UserID:    uuid.MustParse(req.UserID),
				ProductID: uuid.MustParse(product.ProductID),
			}); err != nil {
				logrus.Errorf("[Usecase][UpdateQty] Error on deleting product from cart: %v", err)
				return nil, err
			}
		} else {
			// Update product quantity
			cartItem, err := c.store.UpdateQty(ctx, models.CartRequest{
				UserID:    uuid.MustParse(req.UserID),
				ProductID: uuid.MustParse(product.ProductID),
				Qty:       product.Qty,
			})
			if err != nil {
				logrus.Errorf("[Usecase][UpdateQty] Error on updating product quantity: %v", err)
				return nil, err
			}

			// Convert to GetCartResponse
			cartResponses = append(cartResponses, dto.GetCartResponse{
				ID:        cartItem.ID.String(),
				UserID:    cartItem.UserID.String(),
				ProductID: cartItem.ProductID.String(),
				Qty:       cartItem.Qty,
				CreatedAt: cartItem.CreatedAt,
				UpdatedAt: cartItem.UpdatedAt,
			})
		}
	}

	return cartResponses, nil
}

// DeleteCart is a method that deletes a product from a user's cart.
func (c *cart) DeleteCart(ctx context.Context, req dto.DeleteCartRequest) (string, error) {
	if err := c.store.DeleteProduct(ctx, models.DeleteCartRequest{
		UserID:    uuid.MustParse(req.UserID),
		ProductID: uuid.MustParse(req.ProductID),
	}); err != nil {
		logrus.Errorf("[Usecase][DeleteCart] Error on deleting product from cart: %v", err)
		return "", err
	}

	return "Product deleted from cart", nil
}

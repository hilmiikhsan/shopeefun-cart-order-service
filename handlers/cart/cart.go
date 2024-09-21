package cart

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/dto"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/helpers"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/validators"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/validators/errmsg"
	"github.com/sirupsen/logrus"
)

// cartDto is an interface that defines the methods that our Handler struct depends on.
type cartDto interface {
	GetCartByUserID(ctx context.Context, req dto.GetCartRequest) ([]dto.GetCartResponse, error)
	AddCart(ctx context.Context, req dto.AddCartRequest) (dto.GetCartResponse, error)
	UpdateQty(ctx context.Context, req dto.UpdateCartRequest) ([]dto.GetCartResponse, error)
	DeleteCart(ctx context.Context, req dto.DeleteCartRequest) (string, error)
}

// Handler is a struct that holds a cartDto.
type Handler struct {
	cart      cartDto
	validator *validators.Validator
}

// NewHandler is a constructor function that returns a new Handler.
func NewHandler(cart cartDto, validator *validators.Validator) *Handler {
	return &Handler{
		cart,
		validator,
	}
}

// GetCartByUserID is a handler function to get a cart by user id.
// It first extracts the user id from the URL path, then decodes the request body into a GetCartRequest model.
// It then calls the GetCartByUserID method of the cartDto and sends the helper back to the client.
// GetCartByUserID is a handler function to get a cart by user id.
func (h *Handler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	var req dto.GetCartRequest

	userID := r.PathValue("user_id")
	if userID == "" || userID == ":user_id" {
		logrus.Errorf("[Handler][GetCartByUserID] user_id is required")
		res := helpers.Response{
			Err:    "UserID harus diisi.",
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("[Handler][GetCartByUserID] Failed to decode request body: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	req.UserID = userID

	if err := h.validator.ValidateRequest(req); err != nil {
		logrus.Errorf("[Handler][GetCartByUserID] Validation failed: %v", err)
		code, errorMessages := errmsg.ErrorValidationHandler(err)
		res := helpers.Response{
			Err:    helpers.ErrResponseFieldFormat(errorMessages),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, code)
		return
	}

	data, err := h.cart.GetCartByUserID(r.Context(), req)
	if err != nil {
		logrus.Errorf("[Handler][GetCartByUserID] Failed to get cart by user id: %v", err)
		res := helpers.Response{
			Err:    helpers.STATUS_INTERNAL_ERR,
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusInternalServerError)
		return
	}

	res := helpers.Response{
		Data:   data,
		Msg:    helpers.SUCCESS_RESPONSE,
		Status: true,
	}

	res.HandleResponse(w, http.StatusOK)
}

// AddCart is a handler function to add a product to the cart.
// It decodes the request body into a Cart model, checks if the quantity is greater than 0,
// then calls the AddCart method of the cartDto and sends the helper back to the client.
func (h *Handler) AddCart(w http.ResponseWriter, r *http.Request) {
	var req dto.AddCartRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("[Handler][AddCart] Failed to decode request body: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateRequest(req); err != nil {
		logrus.Errorf("[Handler][AddCart] Validation failed: %v", err)
		code, errorMessages := errmsg.ErrorValidationHandler(err)
		res := helpers.Response{
			Err:    helpers.ErrResponseFieldFormat(errorMessages),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, code)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		logrus.Errorf("[Handler][AddCart] Failed to parse user id: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	if req.Qty <= 0 {
		logrus.Errorf("[Handler][AddCart] Quantity must be greater than 0")
		res := helpers.Response{
			Err:    "Qty harus lebih dari 0.",
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	productResponses, err := h.cart.AddCart(r.Context(), dto.AddCartRequest{
		UserID:    userID.String(),
		ProductID: req.ProductID,
		Qty:       req.Qty,
	})
	if err != nil {
		logrus.Errorf("[Handler][AddCart] Failed to add product to cart: %v", err)

		if strings.Contains(err.Error(), errmsg.ErrFailedToParseProductID) {
			res := helpers.Response{
				Err:    err.Error(),
				Msg:    helpers.FAILED_RESPONSE,
				Status: false,
			}
			res.HandleResponse(w, http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), errmsg.ErrDoesntNewItemsAdded) {
			res := helpers.Response{
				Err:    err.Error(),
				Msg:    helpers.FAILED_RESPONSE,
				Status: false,
			}
			res.HandleResponse(w, http.StatusConflict)
			return
		}

		if strings.Contains(err.Error(), errmsg.ErrUserNotFound) || strings.Contains(err.Error(), errmsg.ErrProductNotFound) {
			res := helpers.Response{
				Err:    err.Error(),
				Msg:    helpers.FAILED_RESPONSE,
				Status: false,
			}
			res.HandleResponse(w, http.StatusNotFound)
			return
		}

		res := helpers.Response{
			Err:    helpers.STATUS_INTERNAL_ERR,
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusInternalServerError)
		return
	}

	res := helpers.Response{
		Data:   productResponses,
		Msg:    helpers.SUCCESS_RESPONSE,
		Status: true,
	}

	res.HandleResponse(w, http.StatusCreated)
}

// UpdateCart is a handler function to delete and update quantity of a product from the cart.
// It first extracts the user id from the URL path, then decodes the request body into a Cart model.
// It then calls the UpdateQty method of the cartDto and sends the helper back to the client.
func (h *Handler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateCartRequest

	userID := r.PathValue("user_id")
	if userID == "" || userID == ":user_id" {
		logrus.Errorf("[Handler][UpdateCart] user_id is required")
		res := helpers.Response{
			Err:    "UserID harus diisi.",
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("[Handler][UpdateCart] Failed to decode request body: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	req.UserID = userID

	if err := h.validator.ValidateNoDuplicateUpdateProductIDs(req.Products); err != nil {
		logrus.Errorf("[Handler][UpdateCart] Validation failed: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateRequest(req); err != nil {
		logrus.Errorf("[Handler][UpdateCart] Validation failed: %v", err)
		code, errorMessages := errmsg.ErrorValidationHandler(err)

		res := helpers.Response{
			Err:    helpers.ErrResponseFieldFormat(errorMessages),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, code)
		return
	}

	response, err := h.cart.UpdateQty(r.Context(), req)
	if err != nil {
		logrus.Errorf("[Handler][UpdateCart] Failed to get cart by user id: %v", err)

		if strings.Contains(err.Error(), errmsg.ErrUserNotFound) || strings.Contains(err.Error(), errmsg.ErrProductNotFound) {
			res := helpers.Response{
				Err:    err.Error(),
				Msg:    helpers.FAILED_RESPONSE,
				Status: false,
			}
			res.HandleResponse(w, http.StatusNotFound)
			return
		}

		res := helpers.Response{
			Err:    helpers.STATUS_INTERNAL_ERR,
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusInternalServerError)
		return
	}

	res := helpers.Response{
		Data:   response,
		Msg:    helpers.SUCCESS_RESPONSE,
		Status: true,
	}

	res.HandleResponse(w, http.StatusOK)
}

// DeleteCart is a handler function to delete a product from the cart.
// It first extracts the user id from the URL path, then decodes the request body into a Cart model.
// It then calls the DeleteCart method of the cartDto and sends the helper back to the client.
func (h *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteCartRequest

	userID := r.PathValue("user_id")
	if userID == "" || userID == ":user_id" {
		logrus.Errorf("[Handler][DeleteCart] user_id is required")
		res := helpers.Response{
			Err:    "UserID harus diisi.",
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("[Handler][DeleteCart] Failed to decode request body: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	req.UserID = userID

	if err := h.validator.ValidateRequest(req); err != nil {
		logrus.Errorf("[Handler][UpdateCart] Validation failed: %v", err)
		code, errorMessages := errmsg.ErrorValidationHandler(err)

		res := helpers.Response{
			Err:    helpers.ErrResponseFieldFormat(errorMessages),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, code)
		return
	}

	response, err := h.cart.DeleteCart(r.Context(), req)
	if err != nil {
		logrus.Errorf("[Handler][DeleteCart] Failed to delete product from cart: %v", err)

		if strings.Contains(err.Error(), errmsg.ErrUserNotFound) || strings.Contains(err.Error(), errmsg.ErrProductNotFound) {
			res := helpers.Response{
				Err:    err.Error(),
				Msg:    helpers.FAILED_RESPONSE,
				Status: false,
			}
			res.HandleResponse(w, http.StatusNotFound)
			return
		}

		res := helpers.Response{
			Err:    helpers.STATUS_INTERNAL_ERR,
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusInternalServerError)
		return
	}

	res := helpers.Response{
		Msg:    response,
		Status: true,
	}

	res.HandleResponse(w, http.StatusOK)
}

package order

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/dto"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/helpers"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/validators"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/validators/errmsg"
	"github.com/sirupsen/logrus"
)

type orderDto interface {
	CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*uuid.UUID, error)
}

type Handler struct {
	order     orderDto
	validator *validators.Validator
}

func NewHandler(order orderDto, validator *validators.Validator) *Handler {
	return &Handler{
		order,
		validator,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("[Handler][CreateOrder] Failed to decode request body: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusBadRequest)
		return
	}

	req.RefCode = helpers.GenerateRefCode()

	if req.ProductOrder == nil {
		logrus.Errorf("[Handler][CreateOrder] Product order is nil")
		req.ProductOrder = json.RawMessage("[]")
	}

	if err := h.validator.ValidateRequest(req); err != nil {
		logrus.Errorf("[Handler][CreateOrder] Validation failed: %v", err)
		code, errorMessages := errmsg.ErrorValidationHandler(err)
		res := helpers.Response{
			Err:    helpers.ErrResponseFieldFormat(errorMessages),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, code)
		return
	}

	result, err := h.order.CreateOrder(r.Context(), req)
	if err != nil {
		logrus.Errorf("[Handler][CreateOrder] Failed to create order: %v", err)
		res := helpers.Response{
			Err:    err.Error(),
			Msg:    helpers.FAILED_RESPONSE,
			Status: false,
		}
		res.HandleResponse(w, http.StatusInternalServerError)
		return
	}

	res := helpers.Response{
		Data:   result,
		Msg:    helpers.SUCCESS_RESPONSE,
		Status: true,
	}

	res.HandleResponse(w, http.StatusCreated)
}

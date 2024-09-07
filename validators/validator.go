package validators

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/hilmiikhsan/shopeefun-cart-order-service/dto"
	"github.com/sirupsen/logrus"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) ValidateRequest(req any) error {
	err := v.validate.Struct(req)
	if err != nil {
		logrus.Errorf("[Validator][ValidateGetCartRequest] Validation failed: %v", err)
		var validationErrors validator.ValidationErrors

		if errors, ok := err.(validator.ValidationErrors); ok {
			validationErrors = errors
		} else {
			logrus.Errorf("[Validator][ValidateGetCartRequest] Failed to cast validation errors: %v", err)
			return fmt.Errorf("validation failed: %v", err)
		}

		return validationErrors
	}

	return nil
}

func (v *Validator) ValidateNoDuplicateUpdateProductIDs(products []dto.ProductRequest) error {
	productIDMap := make(map[string]bool)
	for _, product := range products {
		if productIDMap[product.ProductID] {
			return errors.New("product_id sudah tersedia: " + product.ProductID)
		}
		productIDMap[product.ProductID] = true
	}

	return nil
}

func (v *Validator) ValidateNoDuplicateAddProductIDs(productIDs []string) error {
	productIDMap := make(map[string]bool)
	for _, productID := range productIDs {
		if productIDMap[productID] {
			return errors.New("product_id sudah tersedia: " + productID)
		}
		productIDMap[productID] = true
	}

	return nil
}

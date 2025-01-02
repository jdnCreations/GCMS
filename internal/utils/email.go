package utils

import (
	"github.com/go-playground/validator"
	"github.com/jdnCreations/gcms/internal/models"
)

var validate = validator.New()

func ValidateEmail(customer models.CustomerInfo) error {
	err := validate.Struct(customer)
	if err != nil {
		return err
	}
	return nil
}
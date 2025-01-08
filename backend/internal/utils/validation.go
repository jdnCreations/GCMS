package utils

import (
	"github.com/go-playground/validator"
	"github.com/jdnCreations/gcms/internal/models"
)

var validate = validator.New()

func ValidateCustomer(user interface{}) error {
	err := validate.Struct(user)
	if err != nil {
		return err
	}
	return nil
}

func ValidateGame(game models.GameInfo) error {
	err := validate.Struct(game)
	if err != nil {
		return err
	}
	return nil
}
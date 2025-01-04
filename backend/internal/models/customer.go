package models

type CustomerInfo struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateCustomerInfo struct {
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	Email *string `json:"email"`
}
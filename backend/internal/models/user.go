package models

type UserInfo struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserInfo struct {
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	Email *string `json:"email"`
}
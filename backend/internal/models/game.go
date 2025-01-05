package models

type GameInfo struct {
	Title string `json:"title" validate:"required"`
	Genre string `json:"genre" validate:"required"`
	Copies uint8 `json:"copies" validate:"required"`
}
package models

type GameInfo struct {
	Title string `json:"title" validate:"required"`
	Copies int16 `json:"copies" validate:"required"`
}

type UpdateGameInfo struct {
	Title string `json:"title"`
	Copies *int16 `json:"copies"`
}
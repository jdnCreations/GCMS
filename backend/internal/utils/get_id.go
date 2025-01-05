package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetIdFromRequest(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
  id := vars["id"]
  uuid, err := uuid.Parse(id)
  if err != nil {
		return uuid, err
  }
	return uuid, nil
}
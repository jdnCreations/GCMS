package utils

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetIdFromRequest(idString string, r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
  id := vars[idString]
  uuid, err := uuid.Parse(id)
  if err != nil {
		return uuid, err
  }
	return uuid, nil
}
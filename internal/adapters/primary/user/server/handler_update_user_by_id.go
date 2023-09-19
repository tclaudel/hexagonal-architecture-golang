package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

func (h Handlers) UpdateUserById(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	var userDTO UserUpdate

	err := render.DecodeJSON(r.Body, &userDTO)
	if err != nil {
		JSON(w, r, http.StatusBadRequest, Error{Message: err.Error()})

		return
	}

	user, err := UserUpdateToEntity(id, userDTO)
	if err != nil {
		JSON(w, r, http.StatusBadRequest, Error{Message: err.Error()})

		return
	}

	err = h.user.UpdateUser(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, primary.ErrUserNotFound):
			JSON(w, r, http.StatusNotFound, Error{Message: primary.ErrUserNotFound.Error()})
		default:
			JSON(w, r, http.StatusInternalServerError, Error{Message: primary.ErrUpdatingUser.Error()})
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

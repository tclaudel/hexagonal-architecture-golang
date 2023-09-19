package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

func (h Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDTO User

	err := render.DecodeJSON(r.Body, &userDTO)
	if err != nil {
		JSON(w, r, http.StatusBadRequest, Error{
			Message: err.Error(),
		})

		return
	}

	user, err := UserDTOToEntity(userDTO)
	if err != nil {
		JSON(w, r, http.StatusBadRequest, Error{
			Message: err.Error(),
		})

		return
	}

	err = h.user.CreateUser(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, primary.ErrUserAlreadyExists):
			JSON(w, r, http.StatusConflict, Error{Message: primary.ErrUserAlreadyExists.Error()})
		default:
			JSON(w, r, http.StatusInternalServerError, Error{Message: primary.ErrCreatingUser.Error()})
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}

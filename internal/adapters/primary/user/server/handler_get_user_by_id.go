package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

func (h Handlers) GetUserById(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	user, err := h.user.RetrieveUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, primary.ErrUserNotFound):
			JSON(w, r, http.StatusNotFound, Error{Message: primary.ErrUserNotFound.Error()})
		default:
			JSON(w, r, http.StatusInternalServerError, Error{Message: primary.ErrRetrievingUser.Error()})
		}

		return
	}

	JSON(w, r, http.StatusOK, UserEntityToDTO(user))
}

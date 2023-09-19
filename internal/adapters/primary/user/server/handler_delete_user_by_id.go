package server

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

func (h Handlers) DeleteUserById(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	err := h.user.RemoveUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, primary.ErrUserNotFound):
			JSON(w, r, http.StatusNotFound, Error{Message: primary.ErrUserNotFound.Error()})
		default:
			JSON(w, r, http.StatusInternalServerError, Error{Message: primary.ErrRemovingUser.Error()})
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

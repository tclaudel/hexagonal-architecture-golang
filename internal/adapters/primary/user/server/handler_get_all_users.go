package server

import (
	"net/http"

	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

func (h Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.user.RetrieveUsers(r.Context())
	if err != nil {
		JSON(w, r, http.StatusInternalServerError, Error{
			Message: primary.ErrRetrievingUsers.Error(),
		})

		return
	}

	usersDTO := make(Users, len(users))

	for i, user := range users {
		usersDTO[i] = UserEntityToDTO(user)
	}

	JSON(w, r, http.StatusOK, usersDTO)
}

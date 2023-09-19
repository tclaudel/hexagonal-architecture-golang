package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
)

func (u UserUseCase) RemoveUser(ctx context.Context, id uuid.UUID) error {
	err := u.userRepository.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("error removing user: %s", err.Error())

		switch {
		case errors.Is(err, secondary.ErrUserNotFound):
			return primary.ErrUserNotFound
		default:
			return primary.ErrRemovingUser
		}
	}

	return nil
}

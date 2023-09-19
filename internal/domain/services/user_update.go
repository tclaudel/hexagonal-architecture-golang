package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
)

func (u UserUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	err := u.userRepository.UpdateUser(ctx, user)
	if err != nil {
		log.Printf("Error updating user: %s", err.Error())

		switch {
		case errors.Is(err, secondary.ErrUserNotFound):
			return fmt.Errorf("%w: %s", primary.ErrUserNotFound, err.Error())
		default:
			return fmt.Errorf("%w: %s", primary.ErrUpdatingUser, err.Error())
		}
	}

	return nil
}

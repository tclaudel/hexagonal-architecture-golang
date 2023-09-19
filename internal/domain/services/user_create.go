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

func (u UserUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	err := u.userRepository.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error saving user: %s", err.Error())

		if errors.Is(err, secondary.ErrUserAlreadyExists) {
			return fmt.Errorf("%w: %s", primary.ErrUserAlreadyExists, err.Error())
		}

		return fmt.Errorf("%w: %s", primary.ErrCreatingUser, err.Error())
	}

	return nil
}

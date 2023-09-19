package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
)

func (u UserUseCase) RetrieveUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		log.Printf("error retrieving user: %s", err.Error())

		if errors.Is(err, secondary.ErrUserNotFound) {
			return nil, fmt.Errorf("%w: %s", primary.ErrUserNotFound, err.Error())
		}

		return nil, primary.ErrRetrievingUser
	}

	return user, nil
}

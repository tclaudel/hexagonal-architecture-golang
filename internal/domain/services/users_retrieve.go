package services

import (
	"context"
	"log"

	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
)

func (u UserUseCase) RetrieveUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := u.userRepository.GetUsers(ctx)
	if err != nil {
		log.Printf("error retrieving users: %s", err.Error())

		return nil, primary.ErrRetrievingUsers
	}

	return users, nil
}

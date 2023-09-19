package server

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
)

var (
	ErrInvalidUserDTO = errors.New("error invalid user dto")
)

func UserDTOToEntity(userDTO User) (*entity.User, error) {
	user, err := entity.NewUser(entity.UserParams{
		ID:       userDTO.Id.String(),
		Username: userDTO.Username,
		Email:    userDTO.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserDTO, err.Error())
	}

	return user, nil
}

func UserUpdateToEntity(userID uuid.UUID, userUpdate UserUpdate) (*entity.User, error) {
	user, err := entity.NewUser(entity.UserParams{
		ID:       userID.String(),
		Username: userUpdate.Username,
		Email:    userUpdate.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserDTO, err.Error())
	}

	return user, nil
}

func UserEntityToDTO(user *entity.User) User {
	return User{
		Id:       user.ID(),
		Username: user.Username(),
		Email:    user.Email().Address,
	}
}
